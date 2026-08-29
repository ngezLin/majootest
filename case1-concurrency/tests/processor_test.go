package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"majootest/case1-concurrency/internal/models"
	"majootest/case1-concurrency/internal/processor"
	"majootest/case1-concurrency/internal/reader"
)

func TestValidateAndParse(t *testing.T) {
	tests := []struct {
		name      string
		raw       models.RawRecord
		wantValid bool
		wantError string
	}{
		{
			name: "Valid standard transaction",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 1,
				Columns:    []string{"1", "10", "20", "1500.50", "2026-08-01 12:30:00", "SUCCESS"},
			},
			wantValid: true,
		},
		{
			name: "Invalid column count",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 2,
				Columns:    []string{"1", "10", "1500.50"},
			},
			wantValid: false,
			wantError: "expected 6 columns, got 3",
		},
		{
			name: "Invalid transaction ID",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 3,
				Columns:    []string{"abc", "10", "20", "1500.50", "2026-08-01 12:30:00", "SUCCESS"},
			},
			wantValid: false,
		},
		{
			name: "Invalid bill total",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 4,
				Columns:    []string{"1", "10", "20", "-500", "2026-08-01 12:30:00", "SUCCESS"},
			},
			wantValid: false,
		},
		{
			name: "Invalid timestamp format",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 5,
				Columns:    []string{"1", "10", "20", "500", "01/08/2026", "SUCCESS"},
			},
			wantValid: false,
		},
		{
			name: "Invalid status value",
			raw: models.RawRecord{
				SourceFile: "test.csv",
				LineNumber: 6,
				Columns:    []string{"1", "10", "20", "500", "2026-08-01 12:30:00", "CANCELLED_INVALID"},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, pErr := processor.ValidateAndParse(tt.raw)
			if tt.wantValid {
				if pErr != nil {
					t.Fatalf("expected valid record, got error: %v", pErr)
				}
				if tx == nil || tx.BillTotal != 1500.50 {
					t.Fatalf("unexpected transaction data: %+v", tx)
				}
			} else {
				if pErr == nil {
					t.Fatalf("expected error, got valid transaction")
				}
			}
		})
	}
}

func TestConcurrentPipelineIntegration(t *testing.T) {
	// Create temporary directory with test CSV files
	tmpDir, err := os.MkdirTemp("", "majoo_concurrency_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file1Content := `id,merchant_id,outlet_id,bill_total,transaction_time,status
1,1,1,1000.00,2026-08-01 10:00:00,SUCCESS
2,1,1,2000.00,2026-08-01 11:00:00,SUCCESS
3,1,1,INVALID_AMOUNT,2026-08-01 12:00:00,SUCCESS
`
	file2Content := `id,merchant_id,outlet_id,bill_total,transaction_time,status
4,2,2,3000.00,2026-08-01 13:00:00,PENDING
5,2,2,4000.00,2026-08-01 14:00:00,FAILED
`

	os.WriteFile(filepath.Join(tmpDir, "batch1.csv"), []byte(file1Content), 0644)
	os.WriteFile(filepath.Join(tmpDir, "batch2.csv"), []byte(file2Content), 0644)

	files, err := reader.DiscoverCSVFiles(tmpDir)
	if err != nil {
		t.Fatalf("failed to discover test files: %v", err)
	}

	jobsChan := make(chan models.Job, 100)
	resultsChan := make(chan models.Result, 100)

	csvReader := reader.NewCSVReader(files)
	workerPool := processor.NewWorkerPool(4)
	aggregator := processor.NewAggregator(len(files))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := csvReader.StreamFiles(jobsChan, nil); err != nil {
		t.Fatalf("failed to stream files: %v", err)
	}

	workerPool.Start(ctx, jobsChan, resultsChan, nil)
	aggregator.Consume(resultsChan)

	summary := aggregator.GenerateSummary()

	if summary.TotalRowsRead != 5 {
		t.Errorf("expected 5 rows read, got %d", summary.TotalRowsRead)
	}
	if summary.ValidRows != 4 {
		t.Errorf("expected 4 valid rows, got %d", summary.ValidRows)
	}
	if summary.FailedRows != 1 {
		t.Errorf("expected 1 failed row, got %d", summary.FailedRows)
	}
	expectedRevenue := 1000.00 + 2000.00 + 3000.00 + 4000.00
	if summary.TotalRevenue != expectedRevenue {
		t.Errorf("expected total revenue %.2f, got %.2f", expectedRevenue, summary.TotalRevenue)
	}
}

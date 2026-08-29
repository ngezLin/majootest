package processor

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"majootest/case1-concurrency/internal/models"
)

// Aggregator accumulates stream results into a consolidated summary report.
type Aggregator struct {
	TotalFiles       int
	TotalRowsRead    int64
	ValidRows        int64
	FailedRows       int64
	TotalRevenue     float64
	MinBill          float64
	MaxBill          float64
	RevenueByStatus  map[string]float64
	MerchantRevenue  map[int64]float64
	Errors           []*models.ProcessingError
	StartTime        time.Time
	EndTime          time.Time
}

// NewAggregator creates a new Aggregator instance.
func NewAggregator(totalFiles int) *Aggregator {
	return &Aggregator{
		TotalFiles:      totalFiles,
		MinBill:         math.MaxFloat64,
		MaxBill:         0,
		RevenueByStatus: make(map[string]float64),
		MerchantRevenue: make(map[int64]float64),
		Errors:          make([]*models.ProcessingError, 0),
		StartTime:       time.Now(),
	}
}

// Consume processes all incoming results from resultsChan until the channel closes.
func (a *Aggregator) Consume(resultsChan <-chan models.Result) {
	for res := range resultsChan {
		a.TotalRowsRead++

		if res.Error != nil {
			a.FailedRows++
			a.Errors = append(a.Errors, res.Error)
			continue
		}

		tx := res.Transaction
		a.ValidRows++
		a.TotalRevenue += tx.BillTotal

		if tx.BillTotal < a.MinBill {
			a.MinBill = tx.BillTotal
		}
		if tx.BillTotal > a.MaxBill {
			a.MaxBill = tx.BillTotal
		}

		a.RevenueByStatus[tx.Status] += tx.BillTotal
		a.MerchantRevenue[tx.MerchantID] += tx.BillTotal
	}

	a.EndTime = time.Now()
	if a.ValidRows == 0 {
		a.MinBill = 0
	}
}

// GenerateSummary returns the final SummaryReport struct.
func (a *Aggregator) GenerateSummary() *models.SummaryReport {
	duration := a.EndTime.Sub(a.StartTime)
	durationMs := duration.Milliseconds()
	if durationMs == 0 {
		durationMs = 1
	}

	avgBill := 0.0
	if a.ValidRows > 0 {
		avgBill = a.TotalRevenue / float64(a.ValidRows)
	}

	rps := float64(a.TotalRowsRead) / duration.Seconds()
	if duration.Seconds() == 0 {
		rps = float64(a.TotalRowsRead)
	}

	return &models.SummaryReport{
		TotalFiles:       a.TotalFiles,
		TotalRowsRead:    a.TotalRowsRead,
		ValidRows:        a.ValidRows,
		FailedRows:       a.FailedRows,
		TotalRevenue:     math.Round(a.TotalRevenue*100) / 100,
		AverageBill:      math.Round(avgBill*100) / 100,
		MinBill:          math.Round(a.MinBill*100) / 100,
		MaxBill:          math.Round(a.MaxBill*100) / 100,
		RevenueByStatus:  a.RevenueByStatus,
		MerchantRevenue:  a.MerchantRevenue,
		ProcessingTimeMs: durationMs,
		ThroughputRPS:    math.Round(rps*100) / 100,
	}
}

// ExportSummaryJSON writes the summary report to a JSON file.
func (a *Aggregator) ExportSummaryJSON(outputPath string) error {
	summary := a.GenerateSummary()

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// ExportErrorsLog writes the error audit list to a log file.
func (a *Aggregator) ExportErrorsLog(outputPath string) error {
	if len(a.Errors) == 0 {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create errors log: %w", err)
	}
	defer file.Close()

	for _, pErr := range a.Errors {
		fmt.Fprintf(file, "[%s:Line %d] %s | Raw Data: %v\n",
			pErr.SourceFile, pErr.LineNumber, pErr.Reason, pErr.RawColumns)
	}

	return nil
}

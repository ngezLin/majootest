package models

import "time"

// RawRecord represents an unparsed line extracted from a CSV file.
type RawRecord struct {
	SourceFile string
	LineNumber int
	Columns    []string
}

// Transaction represents a parsed and validated transaction entity.
type Transaction struct {
	ID              int64     `json:"id"`
	MerchantID      int64     `json:"merchant_id"`
	OutletID        int64     `json:"outlet_id"`
	BillTotal       float64   `json:"bill_total"`
	TransactionTime time.Time `json:"transaction_time"`
	Status          string    `json:"status"` // e.g. SUCCESS, PENDING, FAILED
	SourceFile      string    `json:"source_file"`
	LineNumber      int       `json:"line_number"`
}

// Job represents a unit of work sent across the worker channel.
type Job struct {
	Raw RawRecord
}

// Result represents the outcome of processing a Job by a worker.
type Result struct {
	Transaction *Transaction
	Error       *ProcessingError
}

// SummaryReport holds aggregated statistics across all processed CSV files.
type SummaryReport struct {
	TotalFiles       int                `json:"total_files"`
	TotalRowsRead    int64              `json:"total_rows_read"`
	ValidRows        int64              `json:"valid_rows"`
	FailedRows       int64              `json:"failed_rows"`
	TotalRevenue     float64            `json:"total_revenue"`
	AverageBill      float64            `json:"average_bill"`
	MinBill          float64            `json:"min_bill"`
	MaxBill          float64            `json:"max_bill"`
	RevenueByStatus  map[string]float64 `json:"revenue_by_status"`
	MerchantRevenue  map[int64]float64  `json:"merchant_revenue"`
	ProcessingTimeMs int64              `json:"processing_time_ms"`
	ThroughputRPS    float64            `json:"throughput_records_per_sec"`
}

package processor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"majootest/case1-concurrency/internal/models"
)

var allowedTimeFormats = []string{
	"2006-01-02 15:04:05",
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02",
}

// ValidateAndParse processes a raw CSV record and converts it to a validated Transaction.
func ValidateAndParse(raw models.RawRecord) (*models.Transaction, *models.ProcessingError) {
	if len(raw.Columns) != 6 {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("expected 6 columns, got %d", len(raw.Columns)),
		}
	}

	// 1. Parse Transaction ID
	id, err := strconv.ParseInt(strings.TrimSpace(raw.Columns[0]), 10, 64)
	if err != nil || id <= 0 {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid transaction id '%s': must be a positive integer", raw.Columns[0]),
		}
	}

	// 2. Parse Merchant ID
	merchantID, err := strconv.ParseInt(strings.TrimSpace(raw.Columns[1]), 10, 64)
	if err != nil || merchantID <= 0 {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid merchant id '%s': must be a positive integer", raw.Columns[1]),
		}
	}

	// 3. Parse Outlet ID
	outletID, err := strconv.ParseInt(strings.TrimSpace(raw.Columns[2]), 10, 64)
	if err != nil || outletID <= 0 {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid outlet id '%s': must be a positive integer", raw.Columns[2]),
		}
	}

	// 4. Parse Bill Total
	billTotal, err := strconv.ParseFloat(strings.TrimSpace(raw.Columns[3]), 64)
	if err != nil || billTotal < 0 {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid bill total '%s': must be a non-negative number", raw.Columns[3]),
		}
	}

	// 5. Parse Transaction Timestamp
	rawTime := strings.TrimSpace(raw.Columns[4])
	var parsedTime time.Time
	var timeErr error
	for _, layout := range allowedTimeFormats {
		parsedTime, timeErr = time.Parse(layout, rawTime)
		if timeErr == nil {
			break
		}
	}
	if timeErr != nil {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid timestamp format '%s': unsupported date format", raw.Columns[4]),
		}
	}

	// 6. Parse Status
	status := strings.ToUpper(strings.TrimSpace(raw.Columns[5]))
	if status != "SUCCESS" && status != "PENDING" && status != "FAILED" && status != "REFUNDED" {
		return nil, &models.ProcessingError{
			SourceFile: raw.SourceFile,
			LineNumber: raw.LineNumber,
			RawColumns: raw.Columns,
			Reason:     fmt.Sprintf("invalid status '%s': must be SUCCESS, PENDING, FAILED, or REFUNDED", raw.Columns[5]),
		}
	}

	return &models.Transaction{
		ID:              id,
		MerchantID:      merchantID,
		OutletID:        outletID,
		BillTotal:       billTotal,
		TransactionTime: parsedTime,
		Status:          status,
		SourceFile:      raw.SourceFile,
		LineNumber:      raw.LineNumber,
	}, nil
}

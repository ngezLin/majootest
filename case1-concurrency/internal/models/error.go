package models

import "fmt"

// ProcessingError captures detailed information about a failed CSV row.
type ProcessingError struct {
	SourceFile string   `json:"source_file"`
	LineNumber int      `json:"line_number"`
	RawColumns []string `json:"raw_columns"`
	Reason     string   `json:"reason"`
}

func (e *ProcessingError) Error() string {
	return fmt.Sprintf("[%s:%d] parsing error: %s (raw: %v)", e.SourceFile, e.LineNumber, e.Reason, e.RawColumns)
}

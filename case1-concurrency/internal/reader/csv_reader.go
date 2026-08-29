package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"majootest/case1-concurrency/internal/models"
)

// CSVReader coordinates concurrent streaming of CSV files into a jobs channel.
type CSVReader struct {
	filePaths []string
}

// NewCSVReader creates a new CSVReader for a list of filepaths.
func NewCSVReader(filePaths []string) *CSVReader {
	return &CSVReader{filePaths: filePaths}
}

// DiscoverCSVFiles finds all .csv files inside the given directory.
func DiscoverCSVFiles(dirPath string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", dirPath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".csv" {
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .csv files found in directory '%s'", dirPath)
	}

	return files, nil
}

// StreamFiles concurrently streams rows from all CSV files into the jobs channel.
// It closes jobsChan once all files have been completely read.
func (r *CSVReader) StreamFiles(jobsChan chan<- models.Job, totalRowsCounter func(int64)) error {
	var wg sync.WaitGroup

	for _, filePath := range r.filePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			r.readFile(path, jobsChan, totalRowsCounter)
		}(filePath)
	}

	// Wait for all file reading goroutines to complete, then close jobs channel
	go func() {
		wg.Wait()
		close(jobsChan)
	}()

	return nil
}

// readFile opens a single CSV file and streams lines row-by-row into jobsChan.
func (r *CSVReader) readFile(filePath string, jobsChan chan<- models.Job, totalRowsCounter func(int64)) {
	file, err := os.Open(filePath)
	if err != nil {
		// Emit synthetic error job for unreadable file
		jobsChan <- models.Job{
			Raw: models.RawRecord{
				SourceFile: filePath,
				LineNumber: 0,
				Columns:    nil,
			},
		}
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.ReuseRecord = false
	csvReader.FieldsPerRecord = -1 // Allow variable columns to be caught cleanly by validator

	lineNumber := 0
	fileName := filepath.Base(filePath)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		lineNumber++

		// Skip header row if detected on first line
		if lineNumber == 1 && isHeader(record) {
			continue
		}

		if totalRowsCounter != nil {
			totalRowsCounter(1)
		}

		if err != nil {
			// Parsing error from csv reader
			jobsChan <- models.Job{
				Raw: models.RawRecord{
					SourceFile: fileName,
					LineNumber: lineNumber,
					Columns:    []string{fmt.Sprintf("csv read error: %v", err)},
				},
			}
			continue
		}

		jobsChan <- models.Job{
			Raw: models.RawRecord{
				SourceFile: fileName,
				LineNumber: lineNumber,
				Columns:    record,
			},
		}
	}
}

// isHeader checks if the first row is likely a column header.
func isHeader(columns []string) bool {
	if len(columns) == 0 {
		return false
	}
	first := columns[0]
	return first == "id" || first == "ID" || first == "transaction_id" || first == "TransactionID"
}

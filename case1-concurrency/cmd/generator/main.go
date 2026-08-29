package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	outDir := flag.String("out", "./data/sample", "Output directory for generated CSV files")
	numFiles := flag.Int("files", 5, "Number of CSV files to generate")
	rowsPerFile := flag.Int("rows", 10000, "Number of rows per CSV file")
	errorRate := flag.Float64("error-rate", 0.02, "Proportion of intentionally invalid rows (0.0 - 1.0)")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("⚡ Generating %d CSV files with %d rows each in '%s'...\n", *numFiles, *rowsPerFile, *outDir)

	statuses := []string{"SUCCESS", "PENDING", "FAILED", "REFUNDED"}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	globalID := int64(1)

	baseDate := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)

	for fileIdx := 1; fileIdx <= *numFiles; fileIdx++ {
		fileName := fmt.Sprintf("transactions_%02d.csv", fileIdx)
		filePath := filepath.Join(*outDir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file '%s': %v\n", filePath, err)
			continue
		}

		writer := csv.NewWriter(file)
		// Write CSV header
		writer.Write([]string{"id", "merchant_id", "outlet_id", "bill_total", "transaction_time", "status"})

		for r := 0; r < *rowsPerFile; r++ {
			// Introduce occasional deliberate malformed data for resilience testing
			if rng.Float64() < *errorRate {
				switch rng.Intn(4) {
				case 0:
					// Corrupt column count
					writer.Write([]string{fmt.Sprintf("%d", globalID), "1", "100.00"})
				case 1:
					// Invalid bill total
					writer.Write([]string{fmt.Sprintf("%d", globalID), "1", "1", "NOT_A_NUMBER", "2026-08-01 12:00:00", "SUCCESS"})
				case 2:
					// Invalid date format
					writer.Write([]string{fmt.Sprintf("%d", globalID), "1", "1", "5000", "INVALID_DATE_TIME", "SUCCESS"})
				case 3:
					// Invalid status
					writer.Write([]string{fmt.Sprintf("%d", globalID), "1", "1", "5000", "2026-08-01 12:00:00", "UNKNOWN_STATUS"})
				}
				globalID++
				continue
			}

			// Valid row
			merchantID := int64(rng.Intn(10) + 1)
			outletID := int64(rng.Intn(20) + 1)
			billTotal := float64(rng.Intn(50000)+500) / 10.0 // 50.0 to 5000.0
			randomMinutes := rng.Intn(30 * 24 * 60)
			txTime := baseDate.Add(time.Duration(randomMinutes) * time.Minute)
			status := statuses[rng.Intn(len(statuses))]

			writer.Write([]string{
				fmt.Sprintf("%d", globalID),
				fmt.Sprintf("%d", merchantID),
				fmt.Sprintf("%d", outletID),
				fmt.Sprintf("%.2f", billTotal),
				txTime.Format("2006-01-02 15:04:05"),
				status,
			})
			globalID++
		}

		writer.Flush()
		file.Close()
		fmt.Printf("  ✔ Created: %s (%d rows)\n", fileName, *rowsPerFile)
	}

	fmt.Println("✨ Data generation complete!")
}

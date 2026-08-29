package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"majootest/case1-concurrency/internal/models"
	"majootest/case1-concurrency/internal/processor"
	"majootest/case1-concurrency/internal/reader"
	"majootest/case1-concurrency/internal/tracker"
)

func main() {
	inputDir := flag.String("input", "./data/sample", "Path to input directory containing CSV files")
	workers := flag.Int("workers", runtime.NumCPU()*2, "Number of concurrent worker goroutines")
	bufferSize := flag.Int("buffer", 1000, "Buffer size for job and result channels")
	outputJSON := flag.String("output", "./data/output/summary.json", "Path to save output summary JSON")
	errorLog := flag.String("error-log", "./data/output/errors.log", "Path to save error audit log")
	flag.Parse()

	fmt.Println("================================================================")
	fmt.Println("        🚀 Majoo High-Throughput Concurrent CSV Processor        ")
	fmt.Println("================================================================")
	fmt.Printf(" [CONFIG] Input Directory : %s\n", *inputDir)
	fmt.Printf(" [CONFIG] Worker Pool Size: %d workers\n", *workers)
	fmt.Printf(" [CONFIG] Channel Buffer  : %d items\n", *bufferSize)
	fmt.Printf(" [CONFIG] CPU Cores       : %d cores\n", runtime.NumCPU())
	fmt.Println("----------------------------------------------------------------")

	// 1. Discover CSV files
	files, err := reader.DiscoverCSVFiles(*inputDir)
	if err != nil {
		fmt.Printf("❌ Error discovering CSV files: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📂 Discovered %d CSV file(s) for processing.\n\n", len(files))

	// 2. Setup Context with Signal Cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n⚠️  Cancellation requested! Shutting down gracefully...")
		cancel()
	}()

	// 3. Initialize Pipeline Channels
	jobsChan := make(chan models.Job, *bufferSize)
	resultsChan := make(chan models.Result, *bufferSize)

	// 4. Initialize Components
	csvReader := reader.NewCSVReader(files)
	workerPool := processor.NewWorkerPool(*workers)
	aggregator := processor.NewAggregator(len(files))
	progress := tracker.NewProgressTracker()

	// 5. Start Progress Tracker
	progress.Start(80 * time.Millisecond)

	// 6. Start Streaming Reader
	if err := csvReader.StreamFiles(jobsChan, progress.AddTotalRows); err != nil {
		fmt.Printf("❌ Error streaming files: %v\n", err)
		os.Exit(1)
	}

	// 7. Start Worker Pool
	workerPool.Start(ctx, jobsChan, resultsChan, progress.IncrementProcessed)

	// 8. Consume Results in Aggregator (Blocks until resultsChan closes)
	aggregator.Consume(resultsChan)

	// 9. Stop Progress Tracker
	progress.Stop()

	// 10. Generate Final Summary
	summary := aggregator.GenerateSummary()

	// 11. Save JSON Report & Error Logs
	if err := aggregator.ExportSummaryJSON(*outputJSON); err != nil {
		fmt.Printf("⚠️  Failed to save summary JSON: %v\n", err)
	} else {
		fmt.Printf("\n📄 Summary Report saved to: %s\n", *outputJSON)
	}

	if summary.FailedRows > 0 {
		if err := aggregator.ExportErrorsLog(*errorLog); err != nil {
			fmt.Printf("⚠️  Failed to save error log: %v\n", err)
		} else {
			fmt.Printf("⚠️  %d Failed record(s) logged to: %s\n", summary.FailedRows, *errorLog)
		}
	}

	// 12. Print Formatted Report to Console
	printConsoleReport(summary)
}

func printConsoleReport(s *models.SummaryReport) {
	fmt.Println("\n================================================================")
	fmt.Println("                     📊 PROCESSING SUMMARY                      ")
	fmt.Println("================================================================")
	fmt.Printf(" Total Files Processed  : %d files\n", s.TotalFiles)
	fmt.Printf(" Total Rows Read        : %d rows\n", s.TotalRowsRead)
	fmt.Printf(" Valid Transactions     : %d rows (%.2f%%)\n", s.ValidRows, float64(s.ValidRows)/float64(s.TotalRowsRead)*100)
	fmt.Printf(" Failed / Malformed Rows: %d rows (%.2f%%)\n", s.FailedRows, float64(s.FailedRows)/float64(s.TotalRowsRead)*100)
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf(" Total Revenue (Omzet)  : $%.2f\n", s.TotalRevenue)
	fmt.Printf(" Average Bill           : $%.2f\n", s.AverageBill)
	fmt.Printf(" Min / Max Bill Amount  : $%.2f / $%.2f\n", s.MinBill, s.MaxBill)
	fmt.Println("----------------------------------------------------------------")
	fmt.Println(" 📈 Revenue by Status:")
	for status, rev := range s.RevenueByStatus {
		fmt.Printf("    • %-10s : $%.2f\n", status, rev)
	}
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf(" ⏱️  Processing Time     : %d ms\n", s.ProcessingTimeMs)
	fmt.Printf(" ⚡ Overall Throughput  : %.2f records/sec\n", s.ThroughputRPS)
	fmt.Println("================================================================")
}

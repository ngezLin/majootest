package tests

import (
	"context"
	"fmt"
	"testing"

	"majootest/case1-concurrency/internal/models"
	"majootest/case1-concurrency/internal/processor"
)

func BenchmarkWorkerPool(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 16}

	for _, count := range workerCounts {
		b.Run(fmt.Sprintf("Workers-%d", count), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				jobsChan := make(chan models.Job, 1000)
				resultsChan := make(chan models.Result, 1000)

				pool := processor.NewWorkerPool(count)
				aggregator := processor.NewAggregator(1)

				ctx := context.Background()

				// Feed 10,000 synthetic jobs
				go func() {
					for j := 1; j <= 10000; j++ {
						jobsChan <- models.Job{
							Raw: models.RawRecord{
								SourceFile: "bench.csv",
								LineNumber: j,
								Columns:    []string{"1", "1", "1", "100.50", "2026-08-01 12:00:00", "SUCCESS"},
							},
						}
					}
					close(jobsChan)
				}()

				pool.Start(ctx, jobsChan, resultsChan, nil)
				aggregator.Consume(resultsChan)
			}
		})
	}
}

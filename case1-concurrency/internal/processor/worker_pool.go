package processor

import (
	"context"
	"sync"

	"majootest/case1-concurrency/internal/models"
)

// WorkerPool manages concurrent processing of jobs using a fixed set of worker goroutines.
type WorkerPool struct {
	numWorkers int
}

// NewWorkerPool creates a new WorkerPool instance.
func NewWorkerPool(numWorkers int) *WorkerPool {
	if numWorkers <= 0 {
		numWorkers = 4
	}
	return &WorkerPool{numWorkers: numWorkers}
}

// Start launches the worker goroutines and routes results to resultsChan.
// It closes resultsChan when all workers have finished processing.
func (wp *WorkerPool) Start(
	ctx context.Context,
	jobsChan <-chan models.Job,
	resultsChan chan<- models.Result,
	onJobDone func(),
) {
	var wg sync.WaitGroup

	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobsChan:
					if !ok {
						return // jobsChan closed
					}

					tx, pErr := ValidateAndParse(job.Raw)
					resultsChan <- models.Result{
						Transaction: tx,
						Error:       pErr,
					}

					if onJobDone != nil {
						onJobDone()
					}
				}
			}
		}(i + 1)
	}

	// Close resultsChan after all workers finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
}

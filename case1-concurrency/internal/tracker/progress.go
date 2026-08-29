package tracker

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// ProgressTracker provides real-time, thread-safe CLI progress updates.
type ProgressTracker struct {
	totalRows     int64
	processedRows int64
	stopChan      chan struct{}
	startTime     time.Time
}

// NewProgressTracker creates a new tracker instance.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		stopChan:  make(chan struct{}),
		startTime: time.Now(),
	}
}

// SetTotalRows sets the expected total rows (if known in advance).
func (pt *ProgressTracker) SetTotalRows(total int64) {
	atomic.StoreInt64(&pt.totalRows, total)
}

// AddTotalRows increments the total known rows.
func (pt *ProgressTracker) AddTotalRows(delta int64) {
	atomic.AddInt64(&pt.totalRows, delta)
}

// IncrementProcessed increments the processed record counter by 1.
func (pt *ProgressTracker) IncrementProcessed() {
	atomic.AddInt64(&pt.processedRows, 1)
}

// Start launches a background goroutine that refreshes the console progress bar.
func (pt *ProgressTracker) Start(refreshInterval time.Duration) {
	if refreshInterval <= 0 {
		refreshInterval = 100 * time.Millisecond
	}

	ticker := time.NewTicker(refreshInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-pt.stopChan:
				pt.render(true)
				return
			case <-ticker.C:
				pt.render(false)
			}
		}
	}()
}

// Stop signals the progress tracker to stop updating and prints the final state.
func (pt *ProgressTracker) Stop() {
	close(pt.stopChan)
}

// render formats and prints the live progress bar.
func (pt *ProgressTracker) render(isFinal bool) {
	processed := atomic.LoadInt64(&pt.processedRows)
	total := atomic.LoadInt64(&pt.totalRows)
	elapsed := time.Since(pt.startTime).Seconds()

	rps := float64(0)
	if elapsed > 0 {
		rps = float64(processed) / elapsed
	}

	if total > 0 {
		percent := float64(processed) / float64(total) * 100
		if percent > 100 {
			percent = 100
		}
		barWidth := 25
		filled := int(float64(barWidth) * (percent / 100))
		if filled > barWidth {
			filled = barWidth
		}
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

		fmt.Printf("\r\033[K[PROGRESS] [%s] %5.1f%% | %d/%d records | %6.0f rec/s | Elapsed: %.1fs",
			bar, percent, processed, total, rps, elapsed)
	} else {
		fmt.Printf("\r\033[K[PROGRESS] Streaming... | %d records processed | %6.0f rec/s | Elapsed: %.1fs",
			processed, rps, elapsed)
	}

	if isFinal {
		fmt.Println() // Print newline on finish
	}
}

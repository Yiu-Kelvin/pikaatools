package watch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/Yiu-Kelvin/pikaatools/pkg/aws"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

// Watcher handles periodic scanning and comparison
type Watcher struct {
	scanner     *scanner.NetworkScanner
	comparator  *Comparator
	interval    time.Duration
	verbose     bool
	region      string
	vpcID       string
}

// NewWatcher creates a new watcher instance
func NewWatcher(awsClient *aws.Client, interval time.Duration, verbose bool, region, vpcID string) *Watcher {
	return &Watcher{
		scanner:     scanner.NewNetworkScanner(awsClient),
		comparator:  NewComparator(verbose),
		interval:    interval,
		verbose:     verbose,
		region:      region,
		vpcID:       vpcID,
	}
}

// WatchOptions contains options for the watch command
type WatchOptions struct {
	WorkingStateFile string
	Interval         time.Duration
	Region           string
	Profile          string
	VpcID            string
	Verbose          bool
}

// Watch starts watching for changes against a baseline working state
func (w *Watcher) Watch(ctx context.Context, workingStateFile string) error {
	// Load the baseline working state
	if w.verbose {
		fmt.Printf("Loading baseline state from %s...\n", workingStateFile)
	}

	baseline, err := w.comparator.LoadWorkingState(workingStateFile)
	if err != nil {
		return fmt.Errorf("failed to load baseline state: %w", err)
	}

	if w.verbose {
		fmt.Printf("Loaded baseline state from %s (scanned at %s)\n",
			workingStateFile, baseline.ScanTime.Format(time.RFC3339))
		fmt.Printf("Starting periodic scan every %v...\n", w.interval)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a ticker for periodic scans
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Set verbose mode for scanner
	w.scanner.SetVerbose(w.verbose)

	// Perform initial scan
	color.Cyan("🔍 Starting initial scan...")
	if err := w.performScan(ctx, baseline); err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			color.Yellow("Watch stopped by context cancellation")
			return ctx.Err()

		case <-sigChan:
			color.Yellow("\nWatch stopped by signal")
			return nil

		case <-ticker.C:
			color.Cyan("🔍 Performing periodic scan...")
			if err := w.performScan(ctx, baseline); err != nil {
				color.Red("Scan failed: %v", err)
				// Continue watching even if one scan fails
			}
		}
	}
}

// performScan executes a scan and compares against baseline
func (w *Watcher) performScan(ctx context.Context, baseline *scanner.Network) error {
	scanStart := time.Now()

	// Perform the scan
	current, err := w.scanner.ScanNetwork(ctx, w.vpcID)
	if err != nil {
		return fmt.Errorf("failed to scan network: %w", err)
	}

	scanDuration := time.Since(scanStart)

	// Compare with baseline
	differences := w.comparator.Compare(baseline, current)

	// Print timestamp and scan info
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if w.verbose {
		fmt.Printf("\n[%s] Scan completed in %v (region: %s)\n", timestamp, scanDuration, w.region)
	} else {
		fmt.Printf("\n[%s] ", timestamp)
	}

	// Print differences
	w.comparator.PrintDifferences(differences)

	return nil
}
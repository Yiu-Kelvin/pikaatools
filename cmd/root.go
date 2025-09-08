package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/Yiu-Kelvin/pikaatools/pkg/aws"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
	"github.com/Yiu-Kelvin/pikaatools/pkg/graph"
	"github.com/Yiu-Kelvin/pikaatools/pkg/watch"
)

var (
	region       string
	profile      string
	vpcID        string
	output       string
	verbose      bool
	exportJSON   string
	saveState    bool
	
	// Watch command flags
	workingStateFile string
	watchInterval    time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "pikaatools",
	Short: "AWS Network Scanner and Visualizer",
	Long: `PikaaTools is a comprehensive AWS network scanner that discovers and visualizes 
your AWS network infrastructure including VPCs, subnets, peering connections, 
Transit Gateways, IAM roles and policies, and other network resources.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan AWS network infrastructure",
	Long: `Scan your AWS network infrastructure and generate a visual representation
of VPCs, subnets, peering connections, Transit Gateways, IAM roles and policies, and related resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScan(cmd.Context())
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes in AWS network infrastructure",
	Long: `Watch for changes in your AWS network infrastructure by periodically scanning
and comparing against a baseline working state. Displays differences in red text
when changes are detected.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWatch(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(watchCmd)
	
	// Scan command flags
	scanCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region (defaults to AWS_REGION or us-east-1)")
	scanCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile (defaults to default profile)")
	scanCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "Specific VPC ID to scan (scans all VPCs if not provided)")
	scanCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, dot")
	scanCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	scanCmd.Flags().StringVar(&exportJSON, "export-json", "", "Export working state to JSON file (e.g., working_state.json)")
	scanCmd.Flags().BoolVar(&saveState, "save-state", false, "Save working state to working_state.json")
	
	// Watch command flags
	watchCmd.Flags().StringVarP(&workingStateFile, "file", "f", "working_state.json", "Working state file to compare against")
	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 30*time.Second, "Scan interval (e.g., 30s, 1m, 5m)")
	watchCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region (defaults to AWS_REGION or us-east-1)")
	watchCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile (defaults to default profile)")
	watchCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "Specific VPC ID to watch (watches all VPCs if not provided)")
	watchCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func runScan(ctx context.Context) error {
	if verbose {
		fmt.Println("Initializing AWS client...")
	}
	
	// Initialize AWS client
	awsClient, err := aws.NewClient(ctx, region, profile)
	if err != nil {
		return fmt.Errorf("failed to initialize AWS client: %w", err)
	}
	
	if verbose {
		fmt.Printf("Scanning AWS network infrastructure in region: %s\n", awsClient.Region())
	}
	
	// Initialize scanner
	networkScanner := scanner.NewNetworkScanner(awsClient)
	networkScanner.SetVerbose(verbose)
	
	// Scan network infrastructure
	network, err := networkScanner.ScanNetwork(ctx, vpcID)
	if err != nil {
		return fmt.Errorf("failed to scan network: %w", err)
	}
	
	if verbose {
		fmt.Printf("Found %d VPCs, %d subnets, %d peering connections, %d transit gateways, %d security groups, %d network ACLs, %d IAM roles\n", 
			len(network.VPCs), 
			len(network.Subnets),
			len(network.PeeringConnections),
			len(network.TransitGateways),
			len(network.SecurityGroups),
			len(network.NetworkAcls),
			len(network.IAMRoles))
	}
	
	// Set default filename if save-state flag is used
	if saveState && exportJSON == "" {
		exportJSON = "working_state.json"
	}
	
	// Export to JSON if requested
	if exportJSON != "" {
		if verbose {
			fmt.Printf("Exporting working state to %s...\n", exportJSON)
		}
		
		jsonData, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal network data to JSON: %w", err)
		}
		
		err = os.WriteFile(exportJSON, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write JSON file %s: %w", exportJSON, err)
		}
		
		if verbose {
			fmt.Printf("Working state exported successfully to %s\n", exportJSON)
		}
		
		// If only JSON export was requested, don't generate visualization
		if output == "text" && exportJSON != "" {
			return nil
		}
	}
	
	// Generate visualization
	visualizer := graph.NewVisualizer(output)
	result, err := visualizer.Generate(network)
	if err != nil {
		return fmt.Errorf("failed to generate visualization: %w", err)
	}
	
	fmt.Print(result)
	return nil
}

func runWatch(ctx context.Context) error {
	if verbose {
		fmt.Println("Initializing AWS client...")
	}
	
	// Initialize AWS client
	awsClient, err := aws.NewClient(ctx, region, profile)
	if err != nil {
		return fmt.Errorf("failed to initialize AWS client: %w", err)
	}
	
	if verbose {
		fmt.Printf("Starting watch in region: %s with interval: %v\n", awsClient.Region(), watchInterval)
		fmt.Printf("Watching for changes against baseline: %s\n", workingStateFile)
	}
	
	// Check if working state file exists
	if _, err := os.Stat(workingStateFile); os.IsNotExist(err) {
		return fmt.Errorf("working state file %s does not exist. Please run 'scan --save-state' first to create a baseline", workingStateFile)
	}
	
	// Create and start watcher
	watcher := watch.NewWatcher(awsClient, watchInterval, verbose, awsClient.Region(), vpcID)
	
	return watcher.Watch(ctx, workingStateFile)
}
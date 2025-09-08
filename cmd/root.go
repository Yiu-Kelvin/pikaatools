package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Yiu-Kelvin/pikaatools/pkg/aws"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
	"github.com/Yiu-Kelvin/pikaatools/pkg/graph"
)

var (
	region       string
	profile      string
	vpcID        string
	output       string
	verbose      bool
	exportJSON   string
	saveState    bool
)

var rootCmd = &cobra.Command{
	Use:   "pikaatools",
	Short: "AWS Network Scanner and Visualizer",
	Long: `PikaaTools is a comprehensive AWS network scanner that discovers and visualizes 
your AWS network infrastructure including VPCs, subnets, peering connections, 
Transit Gateways, and other network resources.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan AWS network infrastructure",
	Long: `Scan your AWS network infrastructure and generate a visual representation
of VPCs, subnets, peering connections, Transit Gateways, and related resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScan(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	
	scanCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region (defaults to AWS_REGION or us-east-1)")
	scanCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile (defaults to default profile)")
	scanCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "Specific VPC ID to scan (scans all VPCs if not provided)")
	scanCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, dot")
	scanCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	scanCmd.Flags().StringVar(&exportJSON, "export-json", "", "Export working state to JSON file (e.g., working_state.json)")
	scanCmd.Flags().BoolVar(&saveState, "save-state", false, "Save working state to working_state.json")
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
	
	// Scan network infrastructure
	network, err := networkScanner.ScanNetwork(ctx, vpcID)
	if err != nil {
		return fmt.Errorf("failed to scan network: %w", err)
	}
	
	if verbose {
		fmt.Printf("Found %d VPCs, %d subnets, %d peering connections, %d transit gateways, %d security groups\n", 
			len(network.VPCs), 
			len(network.Subnets),
			len(network.PeeringConnections),
			len(network.TransitGateways),
			len(network.SecurityGroups))
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
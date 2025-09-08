package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Yiu-Kelvin/pikaatools/pkg/aws"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
	"github.com/Yiu-Kelvin/pikaatools/pkg/graph"
)

var (
	region    string
	profile   string
	vpcID     string
	output    string
	verbose   bool
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
		fmt.Printf("Found %d VPCs, %d subnets, %d peering connections, %d transit gateways\n", 
			len(network.VPCs), 
			len(network.Subnets),
			len(network.PeeringConnections),
			len(network.TransitGateways))
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
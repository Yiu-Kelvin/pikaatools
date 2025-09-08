package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yiu-Kelvin/pikaatools/pkg/graph"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

func main() {
	// Create a sample network for demonstration
	network := createSampleNetwork()
	
	// Generate text visualization
	fmt.Println("=== TEXT VISUALIZATION ===")
	textVisualizer := graph.NewVisualizer("text")
	textResult, err := textVisualizer.Generate(network)
	if err != nil {
		fmt.Printf("Error generating text visualization: %v\n", err)
		return
	}
	fmt.Print(textResult)
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("=== DOT VISUALIZATION ===")
	
	// Generate DOT visualization
	dotVisualizer := graph.NewVisualizer("dot")
	dotResult, err := dotVisualizer.Generate(network)
	if err != nil {
		fmt.Printf("Error generating DOT visualization: %v\n", err)
		return
	}
	fmt.Print(dotResult)
}

func createSampleNetwork() *scanner.Network {
	now := time.Now()
	
	return &scanner.Network{
		Region:   "us-east-1",
		ScanTime: now,
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345678",
				Name:      "Production VPC",
				CidrBlock: "10.0.0.0/16",
				State:     "available",
				IsDefault: false,
				Tags: map[string]string{
					"Name":        "Production VPC",
					"Environment": "prod",
				},
				Subnets:          []string{"subnet-11111111", "subnet-22222222", "subnet-33333333"},
				InternetGateways: []string{"igw-12345678"},
				NATGateways:      []string{"nat-12345678"},
			},
			{
				ID:        "vpc-87654321",
				Name:      "Development VPC",
				CidrBlock: "172.16.0.0/16",
				State:     "available",
				IsDefault: false,
				Tags: map[string]string{
					"Name":        "Development VPC",
					"Environment": "dev",
				},
				Subnets:          []string{"subnet-44444444"},
				InternetGateways: []string{"igw-87654321"},
			},
			{
				ID:        "vpc-default99",
				Name:      "Default VPC",
				CidrBlock: "172.31.0.0/16",
				State:     "available",
				IsDefault: true,
				Tags: map[string]string{
					"Name": "Default VPC",
				},
				Subnets:          []string{"subnet-55555555"},
				InternetGateways: []string{"igw-default99"},
			},
		},
		Subnets: []scanner.Subnet{
			{
				ID:               "subnet-11111111",
				Name:             "Public Subnet 1",
				VpcID:            "vpc-12345678",
				CidrBlock:        "10.0.1.0/24",
				AvailabilityZone: "us-east-1a",
				State:            "available",
				Type:             "public",
				Tags: map[string]string{
					"Name": "Public Subnet 1",
					"Type": "public",
				},
			},
			{
				ID:               "subnet-22222222",
				Name:             "Private Subnet 1",
				VpcID:            "vpc-12345678",
				CidrBlock:        "10.0.10.0/24",
				AvailabilityZone: "us-east-1a",
				State:            "available",
				Type:             "private",
				Tags: map[string]string{
					"Name": "Private Subnet 1",
					"Type": "private",
				},
			},
			{
				ID:               "subnet-33333333",
				Name:             "Private Subnet 2",
				VpcID:            "vpc-12345678",
				CidrBlock:        "10.0.11.0/24",
				AvailabilityZone: "us-east-1b",
				State:            "available",
				Type:             "private",
				Tags: map[string]string{
					"Name": "Private Subnet 2",
					"Type": "private",
				},
			},
			{
				ID:               "subnet-44444444",
				Name:             "Dev Subnet",
				VpcID:            "vpc-87654321",
				CidrBlock:        "172.16.1.0/24",
				AvailabilityZone: "us-east-1a",
				State:            "available",
				Type:             "public",
				Tags: map[string]string{
					"Name": "Dev Subnet",
					"Type": "public",
				},
			},
			{
				ID:               "subnet-55555555",
				Name:             "Default Subnet",
				VpcID:            "vpc-default99",
				CidrBlock:        "172.31.1.0/20",
				AvailabilityZone: "us-east-1a",
				State:            "available",
				Type:             "public",
				Tags: map[string]string{
					"Name": "Default Subnet",
				},
			},
		},
		PeeringConnections: []scanner.PeeringConnection{
			{
				ID:             "pcx-12345678",
				Name:           "Prod to Dev Peering",
				RequesterVpcID: "vpc-12345678",
				AccepterVpcID:  "vpc-87654321",
				Status:         "active",
				Tags: map[string]string{
					"Name":    "Prod to Dev Peering",
					"Purpose": "cross-env-access",
				},
			},
		},
		TransitGateways: []scanner.TransitGateway{
			{
				ID:    "tgw-12345678",
				Name:  "Main Transit Gateway",
				State: "available",
				Tags: map[string]string{
					"Name":    "Main Transit Gateway",
					"Purpose": "central-routing",
				},
				Attachments: []scanner.TransitGatewayAttachment{
					{
						ID:               "tgw-attach-11111111",
						TransitGatewayID: "tgw-12345678",
						ResourceID:       "vpc-12345678",
						ResourceType:     "vpc",
						State:            "available",
					},
					{
						ID:               "tgw-attach-22222222",
						TransitGatewayID: "tgw-12345678",
						ResourceID:       "vpc-87654321",
						ResourceType:     "vpc",
						State:            "available",
					},
				},
			},
		},
		InternetGateways: []scanner.InternetGateway{
			{
				ID:    "igw-12345678",
				Name:  "Prod IGW",
				VpcID: "vpc-12345678",
				State: "available",
				Tags: map[string]string{
					"Name": "Prod IGW",
				},
			},
			{
				ID:    "igw-87654321",
				Name:  "Dev IGW",
				VpcID: "vpc-87654321",
				State: "available",
				Tags: map[string]string{
					"Name": "Dev IGW",
				},
			},
			{
				ID:    "igw-default99",
				Name:  "Default IGW",
				VpcID: "vpc-default99",
				State: "available",
				Tags: map[string]string{
					"Name": "Default IGW",
				},
			},
		},
		NATGateways: []scanner.NATGateway{
			{
				ID:               "nat-12345678",
				Name:             "Prod NAT Gateway",
				VpcID:            "vpc-12345678",
				SubnetID:         "subnet-11111111",
				State:            "available",
				PublicIP:         "54.123.45.67",
				PrivateIP:        "10.0.1.100",
				ConnectivityType: "public",
				Tags: map[string]string{
					"Name": "Prod NAT Gateway",
				},
			},
		},
	}
}
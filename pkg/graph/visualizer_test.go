package graph

import (
	"strings"
	"testing"
	"time"

	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

func TestNewVisualizer(t *testing.T) {
	v := NewVisualizer("text")
	if v.format != "text" {
		t.Errorf("Expected format 'text', got '%s'", v.format)
	}
	
	v = NewVisualizer("dot")
	if v.format != "dot" {
		t.Errorf("Expected format 'dot', got '%s'", v.format)
	}
}

func TestGenerateUnsupportedFormat(t *testing.T) {
	v := NewVisualizer("unsupported")
	network := &scanner.Network{}
	
	_, err := v.Generate(network)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
	
	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("Expected 'unsupported output format' error, got: %s", err.Error())
	}
}

func TestGenerateTextGraph(t *testing.T) {
	v := NewVisualizer("text")
	
	network := &scanner.Network{
		Region:   "us-east-1",
		ScanTime: time.Now(),
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "Test VPC",
				CidrBlock: "10.0.0.0/16",
				IsDefault: false,
				State:     "available",
				Tags:      map[string]string{"Name": "Test VPC"},
			},
		},
		Subnets: []scanner.Subnet{
			{
				ID:               "subnet-12345",
				Name:             "Test Subnet",
				VpcID:            "vpc-12345",
				CidrBlock:        "10.0.1.0/24",
				AvailabilityZone: "us-east-1a",
				State:            "available",
				Type:             "public",
				Tags:             map[string]string{"Name": "Test Subnet"},
			},
		},
	}
	
	// Update VPC associations
	network.VPCs[0].Subnets = []string{"subnet-12345"}
	
	result, err := v.Generate(network)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	
	// Check that the output contains expected elements
	if !strings.Contains(result, "Test VPC") {
		t.Error("Expected output to contain 'Test VPC'")
	}
	
	if !strings.Contains(result, "Test Subnet") {
		t.Error("Expected output to contain 'Test Subnet'")
	}
	
	if !strings.Contains(result, "10.0.0.0/16") {
		t.Error("Expected output to contain VPC CIDR")
	}
	
	if !strings.Contains(result, "10.0.1.0/24") {
		t.Error("Expected output to contain subnet CIDR")
	}
	
	if !strings.Contains(result, "Summary:") {
		t.Error("Expected output to contain summary")
	}
}

func TestGenerateDotGraph(t *testing.T) {
	v := NewVisualizer("dot")
	
	network := &scanner.Network{
		Region:   "us-east-1",
		ScanTime: time.Now(),
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "Test VPC",
				CidrBlock: "10.0.0.0/16",
				IsDefault: false,
				State:     "available",
			},
		},
		PeeringConnections: []scanner.PeeringConnection{
			{
				ID:             "pcx-12345",
				Name:           "Test Peering",
				RequesterVpcID: "vpc-12345",
				AccepterVpcID:  "vpc-67890",
				Status:         "active",
			},
		},
	}
	
	result, err := v.Generate(network)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	
	// Check DOT format structure
	if !strings.Contains(result, "digraph AWSNetwork") {
		t.Error("Expected DOT graph to contain 'digraph AWSNetwork'")
	}
	
	if !strings.Contains(result, "vpc-12345") {
		t.Error("Expected DOT graph to contain VPC ID")
	}
	
	if !strings.Contains(result, "Test VPC") {
		t.Error("Expected DOT graph to contain VPC name")
	}
	
	if !strings.Contains(result, "Test Peering") {
		t.Error("Expected DOT graph to contain peering connection")
	}
	
	// Check that it ends properly
	if !strings.HasSuffix(strings.TrimSpace(result), "}") {
		t.Error("Expected DOT graph to end with '}'")
	}
}
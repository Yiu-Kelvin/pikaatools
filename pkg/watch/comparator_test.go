package watch

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

func TestComparator(t *testing.T) {
	comparator := NewComparator(false)
	
	if comparator == nil {
		t.Error("Expected non-nil comparator")
	}
	
	if comparator.verbose {
		t.Error("Expected verbose to be false")
	}
}

func TestLoadWorkingState(t *testing.T) {
	// Create a temporary working state file
	network := &scanner.Network{
		Region:   "us-east-1",
		ScanTime: time.Now(),
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
		},
	}
	
	data, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	
	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "test_working_state_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	// Test loading
	comparator := NewComparator(false)
	loaded, err := comparator.LoadWorkingState(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load working state: %v", err)
	}
	
	if loaded.Region != "us-east-1" {
		t.Errorf("Expected region us-east-1, got %s", loaded.Region)
	}
	
	if len(loaded.VPCs) != 1 {
		t.Errorf("Expected 1 VPC, got %d", len(loaded.VPCs))
	}
	
	if loaded.VPCs[0].ID != "vpc-12345" {
		t.Errorf("Expected VPC ID vpc-12345, got %s", loaded.VPCs[0].ID)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	comparator := NewComparator(false)
	_, err := comparator.LoadWorkingState("non_existent_file.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestCompareIdenticalNetworks(t *testing.T) {
	network := &scanner.Network{
		Region:   "us-east-1",
		ScanTime: time.Now(),
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
		},
	}
	
	comparator := NewComparator(false)
	differences := comparator.Compare(network, network)
	
	if len(differences) != 0 {
		t.Errorf("Expected no differences for identical networks, got %d", len(differences))
	}
}

func TestCompareNetworksWithNewVPC(t *testing.T) {
	baseline := &scanner.Network{
		Region: "us-east-1",
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
		},
	}
	
	current := &scanner.Network{
		Region: "us-east-1",
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
			{
				ID:        "vpc-67890",
				Name:      "new-vpc",
				CidrBlock: "10.1.0.0/16",
			},
		},
	}
	
	comparator := NewComparator(false)
	differences := comparator.Compare(baseline, current)
	
	if len(differences) != 1 {
		t.Errorf("Expected 1 difference, got %d", len(differences))
	}
	
	if differences[0].Type != Added {
		t.Errorf("Expected Added difference type, got %v", differences[0].Type)
	}
	
	if differences[0].ResourceType != "VPC" {
		t.Errorf("Expected VPC resource type, got %s", differences[0].ResourceType)
	}
	
	if differences[0].ResourceID != "vpc-67890" {
		t.Errorf("Expected vpc-67890 resource ID, got %s", differences[0].ResourceID)
	}
}

func TestCompareNetworksWithRemovedVPC(t *testing.T) {
	baseline := &scanner.Network{
		Region: "us-east-1",
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
			{
				ID:        "vpc-67890",
				Name:      "old-vpc",
				CidrBlock: "10.1.0.0/16",
			},
		},
	}
	
	current := &scanner.Network{
		Region: "us-east-1",
		VPCs: []scanner.VPC{
			{
				ID:        "vpc-12345",
				Name:      "test-vpc",
				CidrBlock: "10.0.0.0/16",
			},
		},
	}
	
	comparator := NewComparator(false)
	differences := comparator.Compare(baseline, current)
	
	if len(differences) != 1 {
		t.Errorf("Expected 1 difference, got %d", len(differences))
	}
	
	if differences[0].Type != Removed {
		t.Errorf("Expected Removed difference type, got %v", differences[0].Type)
	}
	
	if differences[0].ResourceID != "vpc-67890" {
		t.Errorf("Expected vpc-67890 resource ID, got %s", differences[0].ResourceID)
	}
}

func TestCompareNetworkAcls(t *testing.T) {
	baseline := &scanner.Network{
		Region: "us-east-1",
		NetworkAcls: []scanner.NetworkAcl{
			{
				ID:    "acl-12345",
				VpcID: "vpc-12345",
				Name:  "test-nacl",
			},
		},
	}
	
	current := &scanner.Network{
		Region: "us-east-1",
		NetworkAcls: []scanner.NetworkAcl{
			{
				ID:    "acl-12345",
				VpcID: "vpc-12345",
				Name:  "test-nacl",
			},
			{
				ID:    "acl-67890",
				VpcID: "vpc-12345",
				Name:  "new-nacl",
			},
		},
	}
	
	comparator := NewComparator(false)
	differences := comparator.Compare(baseline, current)
	
	if len(differences) != 1 {
		t.Errorf("Expected 1 difference, got %d", len(differences))
	}
	
	if differences[0].Type != Added {
		t.Errorf("Expected Added difference type, got %v", differences[0].Type)
	}
	
	if differences[0].ResourceType != "NetworkACL" {
		t.Errorf("Expected NetworkACL resource type, got %s", differences[0].ResourceType)
	}
}

func TestDifferenceTypes(t *testing.T) {
	tests := []struct {
		name     string
		diffType DifferenceType
		expected string
	}{
		{"Added", Added, "Added"},
		{"Removed", Removed, "Removed"},
		{"Modified", Modified, "Modified"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that the types exist and can be used
			diff := Difference{
				Type:         tt.diffType,
				ResourceType: "Test",
				ResourceID:   "test-id",
				Description:  "test description",
			}
			
			if diff.Type != tt.diffType {
				t.Errorf("Expected type %v, got %v", tt.diffType, diff.Type)
			}
		})
	}
}

func TestShouldSkipField(t *testing.T) {
	comparator := NewComparator(false)
	
	skipFields := []string{"ScanTime", "CreateDate", "UpdateDate"}
	normalFields := []string{"ID", "Name", "VpcID", "CidrBlock"}
	
	for _, field := range skipFields {
		if !comparator.shouldSkipField(field) {
			t.Errorf("Expected field %s to be skipped", field)
		}
	}
	
	for _, field := range normalFields {
		if comparator.shouldSkipField(field) {
			t.Errorf("Expected field %s not to be skipped", field)
		}
	}
}
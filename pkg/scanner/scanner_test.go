package scanner

import (
	"testing"
	"time"
)

func TestConvertTags(t *testing.T) {
	// This test doesn't require AWS credentials as it tests a pure function
	
	// Test empty tags
	tags := convertTags(nil)
	if len(tags) != 0 {
		t.Errorf("Expected empty tags map, got %d items", len(tags))
	}
	
	// Test normal case would require AWS SDK types, so we'll keep it simple
	// This demonstrates the testing structure for when we have more complex logic
}

func TestDetermineSubnetType(t *testing.T) {
	tests := []struct {
		name     string
		routes   []Route
		igws     []InternetGateway
		expected string
	}{
		{
			name: "Public subnet with IGW route",
			routes: []Route{
				{
					DestinationCidr: "0.0.0.0/0",
					GatewayID:       "igw-12345",
					State:           "active",
				},
			},
			igws: []InternetGateway{
				{
					ID:    "igw-12345",
					State: "available",
				},
			},
			expected: "public",
		},
		{
			name: "Private subnet with NAT route",
			routes: []Route{
				{
					DestinationCidr: "0.0.0.0/0",
					GatewayID:       "nat-12345",
					State:           "active",
				},
			},
			igws:     []InternetGateway{},
			expected: "private",
		},
		{
			name: "Isolated subnet",
			routes: []Route{
				{
					DestinationCidr: "10.0.0.0/16",
					GatewayID:       "local",
					State:           "active",
				},
			},
			igws:     []InternetGateway{},
			expected: "isolated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeTable := &RouteTable{
				Routes: tt.routes,
			}
			
			result := determineSubnetType(routeTable, tt.igws)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNetworkStructure(t *testing.T) {
	// Test basic network structure
	network := &Network{
		ScanTime: time.Now(),
		Region:   "us-east-1",
	}
	
	if network.ScanTime.IsZero() {
		t.Error("Expected non-zero scan time")
	}
	
	if network.Region != "us-east-1" {
		t.Errorf("Expected region us-east-1, got %s", network.Region)
	}
}

func TestIAMStructure(t *testing.T) {
	// Test IAM role structure
	role := IAMRole{
		ID:                   "AROA123456789",
		Name:                 "test-role",
		Path:                 "/",
		Arn:                  "arn:aws:iam::123456789012:role/test-role",
		Description:          "Test role",
		CreateDate:           time.Now(),
		AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
		MaxSessionDuration:   3600,
		Tags:                 map[string]string{"Environment": "test"},
		AttachedPolicies:     []IAMPolicy{},
		InlinePolicies:       []IAMInlinePolicy{},
	}
	
	if role.Name != "test-role" {
		t.Errorf("Expected role name 'test-role', got %s", role.Name)
	}
	
	if role.MaxSessionDuration != 3600 {
		t.Errorf("Expected max session duration 3600, got %d", role.MaxSessionDuration)
	}
	
	if role.Tags["Environment"] != "test" {
		t.Error("Expected Environment tag to be 'test'")
	}
}

func TestConvertIAMTags(t *testing.T) {
	// Test convertIAMTags function
	tags := convertIAMTags(nil)
	if len(tags) != 0 {
		t.Errorf("Expected empty tags map, got %d items", len(tags))
	}
}

func TestSecurityGroupRuleStructure(t *testing.T) {
	// Test SecurityGroupRule structure
	rule := SecurityGroupRule{
		IpProtocol:                 "tcp",
		FromPort:                   80,
		ToPort:                     80,
		CidrBlocks:                 []string{"0.0.0.0/0"},
		Ipv6CidrBlocks:             []string{"::/0"},
		PrefixListIds:              []string{"pl-12345"},
		ReferencedGroupId:          "sg-12345",
		ReferencedGroupOwnerId:     "123456789012",
		Description:                "Allow HTTP traffic",
		Tags:                       map[string]string{"Name": "HTTP rule"},
	}
	
	if rule.IpProtocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got %s", rule.IpProtocol)
	}
	
	if rule.FromPort != 80 {
		t.Errorf("Expected from port 80, got %d", rule.FromPort)
	}
	
	if rule.ToPort != 80 {
		t.Errorf("Expected to port 80, got %d", rule.ToPort)
	}
	
	if len(rule.CidrBlocks) != 1 || rule.CidrBlocks[0] != "0.0.0.0/0" {
		t.Error("Expected CIDR block '0.0.0.0/0'")
	}
	
	if rule.Description != "Allow HTTP traffic" {
		t.Errorf("Expected description 'Allow HTTP traffic', got %s", rule.Description)
	}
}

func TestSecurityGroupWithRules(t *testing.T) {
	// Test SecurityGroup with rules
	sg := SecurityGroup{
		ID:          "sg-12345",
		Name:        "test-sg",
		Description: "Test security group",
		VpcID:       "vpc-12345",
		Tags:        map[string]string{"Name": "test-sg"},
		IngressRules: []SecurityGroupRule{
			{
				IpProtocol: "tcp",
				FromPort:   80,
				ToPort:     80,
				CidrBlocks: []string{"0.0.0.0/0"},
			},
		},
		EgressRules: []SecurityGroupRule{
			{
				IpProtocol: "tcp",
				FromPort:   443,
				ToPort:     443,
				CidrBlocks: []string{"0.0.0.0/0"},
			},
		},
	}
	
	if sg.ID != "sg-12345" {
		t.Errorf("Expected SG ID 'sg-12345', got %s", sg.ID)
	}
	
	if len(sg.IngressRules) != 1 {
		t.Errorf("Expected 1 ingress rule, got %d", len(sg.IngressRules))
	}
	
	if len(sg.EgressRules) != 1 {
		t.Errorf("Expected 1 egress rule, got %d", len(sg.EgressRules))
	}
	
	if sg.IngressRules[0].FromPort != 80 {
		t.Errorf("Expected ingress rule port 80, got %d", sg.IngressRules[0].FromPort)
	}
	
	if sg.EgressRules[0].FromPort != 443 {
		t.Errorf("Expected egress rule port 443, got %d", sg.EgressRules[0].FromPort)
	}
}

func TestNetworkScannerVerbose(t *testing.T) {
	// Test that NetworkScanner can toggle verbose mode
	scanner := &NetworkScanner{
		client:  nil, // Not testing actual scanning, just the verbose flag
		verbose: false,
	}
	
	if scanner.verbose {
		t.Error("Expected verbose to be false by default")
	}
	
	scanner.SetVerbose(true)
	if !scanner.verbose {
		t.Error("Expected verbose to be true after setting")
	}
	
	scanner.SetVerbose(false)
	if scanner.verbose {
		t.Error("Expected verbose to be false after setting")
	}
}

func TestNetworkAclStructure(t *testing.T) {
	// Test NetworkAcl structure
	nacl := NetworkAcl{
		ID:        "acl-12345",
		Name:      "test-nacl",
		VpcID:     "vpc-12345",
		IsDefault: false,
		Tags:      map[string]string{"Environment": "test"},
		Entries: []NetworkAclEntry{
			{
				RuleNumber: 100,
				Protocol:   "tcp",
				RuleAction: "allow",
				CidrBlock:  "10.0.0.0/16",
				PortRange: &NetworkAclPortRange{
					From: 80,
					To:   80,
				},
				Egress: false,
			},
		},
		Associations: []string{"subnet-12345"},
	}
	
	if nacl.ID != "acl-12345" {
		t.Errorf("Expected NACL ID 'acl-12345', got %s", nacl.ID)
	}
	
	if nacl.IsDefault {
		t.Error("Expected IsDefault to be false")
	}
	
	if len(nacl.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(nacl.Entries))
	}
	
	entry := nacl.Entries[0]
	if entry.RuleNumber != 100 {
		t.Errorf("Expected rule number 100, got %d", entry.RuleNumber)
	}
	
	if entry.Protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got %s", entry.Protocol)
	}
	
	if entry.RuleAction != "allow" {
		t.Errorf("Expected rule action 'allow', got %s", entry.RuleAction)
	}
	
	if entry.PortRange == nil {
		t.Error("Expected port range to be set")
	} else if entry.PortRange.From != 80 || entry.PortRange.To != 80 {
		t.Errorf("Expected port range 80-80, got %d-%d", entry.PortRange.From, entry.PortRange.To)
	}
	
	if entry.Egress {
		t.Error("Expected egress to be false")
	}
	
	if len(nacl.Associations) != 1 || nacl.Associations[0] != "subnet-12345" {
		t.Error("Expected association with subnet-12345")
	}
}

func TestNetworkAclEntryWithIcmp(t *testing.T) {
	// Test NetworkAclEntry with ICMP type
	entry := NetworkAclEntry{
		RuleNumber: 200,
		Protocol:   "icmp",
		RuleAction: "allow",
		CidrBlock:  "0.0.0.0/0",
		IcmpType: &NetworkAclIcmpType{
			Type: 8,
			Code: 0,
		},
		Egress: true,
	}
	
	if entry.Protocol != "icmp" {
		t.Errorf("Expected protocol 'icmp', got %s", entry.Protocol)
	}
	
	if entry.IcmpType == nil {
		t.Error("Expected ICMP type to be set")
	} else if entry.IcmpType.Type != 8 || entry.IcmpType.Code != 0 {
		t.Errorf("Expected ICMP type 8 code 0, got type %d code %d", entry.IcmpType.Type, entry.IcmpType.Code)
	}
	
	if !entry.Egress {
		t.Error("Expected egress to be true")
	}
}

func TestNetworkWithNacls(t *testing.T) {
	// Test Network structure includes NetworkAcls
	network := &Network{
		NetworkAcls: []NetworkAcl{
			{
				ID:    "acl-12345",
				VpcID: "vpc-12345",
			},
		},
		ScanTime: time.Now(),
		Region:   "us-east-1",
	}
	
	if len(network.NetworkAcls) != 1 {
		t.Errorf("Expected 1 Network ACL, got %d", len(network.NetworkAcls))
	}
	
	if network.NetworkAcls[0].ID != "acl-12345" {
		t.Errorf("Expected Network ACL ID 'acl-12345', got %s", network.NetworkAcls[0].ID)
	}
}
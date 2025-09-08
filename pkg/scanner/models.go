package scanner

import (
	"time"
)

// Network represents the complete AWS network infrastructure
type Network struct {
	VPCs                []VPC                 `json:"vpcs"`
	Subnets             []Subnet              `json:"subnets"`
	PeeringConnections  []PeeringConnection   `json:"peering_connections"`
	TransitGateways     []TransitGateway      `json:"transit_gateways"`
	InternetGateways    []InternetGateway     `json:"internet_gateways"`
	NATGateways         []NATGateway          `json:"nat_gateways"`
	RouteTables         []RouteTable          `json:"route_tables"`
	SecurityGroups      []SecurityGroup       `json:"security_groups"`
	NetworkAcls         []NetworkAcl          `json:"network_acls"`
	IAMRoles            []IAMRole             `json:"iam_roles"`
	ScanTime            time.Time             `json:"scan_time"`
	Region              string                `json:"region"`
}

// VPC represents an AWS VPC
type VPC struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	CidrBlock         string            `json:"cidr_block"`
	State             string            `json:"state"`
	IsDefault         bool              `json:"is_default"`
	DhcpOptionsID     string            `json:"dhcp_options_id"`
	Tags              map[string]string `json:"tags"`
	Subnets           []string          `json:"subnets"`           // Subnet IDs
	SecurityGroups    []string          `json:"security_groups"`    // Security Group IDs
	InternetGateways  []string          `json:"internet_gateways"`  // Internet Gateway IDs
	NATGateways       []string          `json:"nat_gateways"`       // NAT Gateway IDs
	NetworkAcls       []string          `json:"network_acls"`       // Network ACL IDs
}

// Subnet represents an AWS subnet
type Subnet struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	VpcID             string            `json:"vpc_id"`
	CidrBlock         string            `json:"cidr_block"`
	AvailabilityZone  string            `json:"availability_zone"`
	State             string            `json:"state"`
	MapPublicIP       bool              `json:"map_public_ip"`
	Tags              map[string]string `json:"tags"`
	RouteTableID      string            `json:"route_table_id"`
	NetworkAclID      string            `json:"network_acl_id"`
	Type              string            `json:"type"` // "public", "private", "isolated"
}

// PeeringConnection represents a VPC peering connection
type PeeringConnection struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	RequesterVpcID   string            `json:"requester_vpc_id"`
	AccepterVpcID    string            `json:"accepter_vpc_id"`
	Status           string            `json:"status"`
	Tags             map[string]string `json:"tags"`
}

// TransitGateway represents an AWS Transit Gateway
type TransitGateway struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	State       string                     `json:"state"`
	Tags        map[string]string          `json:"tags"`
	Attachments []TransitGatewayAttachment `json:"attachments"`
}

// TransitGatewayAttachment represents a TGW attachment
type TransitGatewayAttachment struct {
	ID                 string            `json:"id"`
	TransitGatewayID   string            `json:"transit_gateway_id"`
	ResourceID         string            `json:"resource_id"`
	ResourceType       string            `json:"resource_type"`
	State              string            `json:"state"`
	Tags               map[string]string `json:"tags"`
}

// InternetGateway represents an AWS Internet Gateway
type InternetGateway struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	VpcID string            `json:"vpc_id"`
	State string            `json:"state"`
	Tags  map[string]string `json:"tags"`
}

// NATGateway represents an AWS NAT Gateway
type NATGateway struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	VpcID            string            `json:"vpc_id"`
	SubnetID         string            `json:"subnet_id"`
	State            string            `json:"state"`
	PublicIP         string            `json:"public_ip"`
	PrivateIP        string            `json:"private_ip"`
	ConnectivityType string            `json:"connectivity_type"`
	Tags             map[string]string `json:"tags"`
}

// RouteTable represents an AWS route table
type RouteTable struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	VpcID        string            `json:"vpc_id"`
	IsMain       bool              `json:"is_main"`
	Tags         map[string]string `json:"tags"`
	Routes       []Route           `json:"routes"`
	Associations []string          `json:"associations"` // Subnet IDs
}

// Route represents a route in a route table
type Route struct {
	DestinationCidr    string `json:"destination_cidr"`
	GatewayID          string `json:"gateway_id"`
	InstanceID         string `json:"instance_id"`
	NetworkInterfaceID string `json:"network_interface_id"`
	VpcPeeringID       string `json:"vpc_peering_id"`
	TransitGatewayID   string `json:"transit_gateway_id"`
	State              string `json:"state"`
	Origin             string `json:"origin"`
}

// SecurityGroup represents an AWS security group
type SecurityGroup struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	VpcID        string                `json:"vpc_id"`
	Tags         map[string]string     `json:"tags"`
	IngressRules []SecurityGroupRule   `json:"ingress_rules"`
	EgressRules  []SecurityGroupRule   `json:"egress_rules"`
}

// SecurityGroupRule represents an AWS security group rule
type SecurityGroupRule struct {
	IpProtocol                 string            `json:"ip_protocol"`
	FromPort                   int32             `json:"from_port"`
	ToPort                     int32             `json:"to_port"`
	CidrBlocks                 []string          `json:"cidr_blocks"`
	Ipv6CidrBlocks             []string          `json:"ipv6_cidr_blocks"`
	PrefixListIds              []string          `json:"prefix_list_ids"`
	ReferencedGroupId          string            `json:"referenced_group_id"`
	ReferencedGroupOwnerId     string            `json:"referenced_group_owner_id"`
	Description                string            `json:"description"`
	Tags                       map[string]string `json:"tags"`
}

// IAMRole represents an AWS IAM role
type IAMRole struct {
	ID                   string              `json:"id"`
	Name                 string              `json:"name"`
	Path                 string              `json:"path"`
	Arn                  string              `json:"arn"`
	Description          string              `json:"description"`
	CreateDate           time.Time           `json:"create_date"`
	AssumeRolePolicyDocument string         `json:"assume_role_policy_document"`
	MaxSessionDuration   int32               `json:"max_session_duration"`
	Tags                 map[string]string   `json:"tags"`
	AttachedPolicies     []IAMPolicy         `json:"attached_policies"`
	InlinePolicies       []IAMInlinePolicy   `json:"inline_policies"`
}

// IAMPolicy represents an AWS IAM policy (managed policy)
type IAMPolicy struct {
	Arn                    string            `json:"arn"`
	PolicyName             string            `json:"policy_name"`
	PolicyId               string            `json:"policy_id"`
	Path                   string            `json:"path"`
	DefaultVersionId       string            `json:"default_version_id"`
	AttachmentCount        int32             `json:"attachment_count"`
	PermissionsBoundaryUsageCount int32     `json:"permissions_boundary_usage_count"`
	IsAttachable           bool              `json:"is_attachable"`
	Description            string            `json:"description"`
	CreateDate             time.Time         `json:"create_date"`
	UpdateDate             time.Time         `json:"update_date"`
	Tags                   map[string]string `json:"tags"`
	PolicyDocument         string            `json:"policy_document"`
}

// IAMInlinePolicy represents an inline policy attached to a role
type IAMInlinePolicy struct {
	PolicyName     string `json:"policy_name"`
	PolicyDocument string `json:"policy_document"`
}

// NetworkAcl represents an AWS Network ACL
type NetworkAcl struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	VpcID        string            `json:"vpc_id"`
	IsDefault    bool              `json:"is_default"`
	Tags         map[string]string `json:"tags"`
	Entries      []NetworkAclEntry `json:"entries"`
	Associations []string          `json:"associations"` // Subnet IDs
}

// NetworkAclEntry represents an entry in a Network ACL
type NetworkAclEntry struct {
	RuleNumber   int32  `json:"rule_number"`
	Protocol     string `json:"protocol"`
	RuleAction   string `json:"rule_action"`
	CidrBlock    string `json:"cidr_block"`
	Ipv6CidrBlock string `json:"ipv6_cidr_block"`
	PortRange    *NetworkAclPortRange `json:"port_range,omitempty"`
	IcmpType     *NetworkAclIcmpType  `json:"icmp_type,omitempty"`
	Egress       bool   `json:"egress"`
}

// NetworkAclPortRange represents a port range in a Network ACL entry
type NetworkAclPortRange struct {
	From int32 `json:"from"`
	To   int32 `json:"to"`
}

// NetworkAclIcmpType represents ICMP type and code in a Network ACL entry
type NetworkAclIcmpType struct {
	Type int32 `json:"type"`
	Code int32 `json:"code"`
}
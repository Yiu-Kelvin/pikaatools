package scanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/Yiu-Kelvin/pikaatools/pkg/aws"
)

// NetworkScanner scans AWS network infrastructure
type NetworkScanner struct {
	client *aws.Client
}

// NewNetworkScanner creates a new network scanner
func NewNetworkScanner(client *aws.Client) *NetworkScanner {
	return &NetworkScanner{
		client: client,
	}
}

// ScanNetwork scans the complete network infrastructure
func (s *NetworkScanner) ScanNetwork(ctx context.Context, vpcID string) (*Network, error) {
	network := &Network{
		ScanTime: time.Now(),
		Region:   s.client.Region(),
	}

	// Scan VPCs
	vpcs, err := s.scanVPCs(ctx, vpcID)
	if err != nil {
		return nil, fmt.Errorf("failed to scan VPCs: %w", err)
	}
	network.VPCs = vpcs

	// Get VPC IDs for filtering other resources
	vpcIDs := make([]string, len(vpcs))
	for i, vpc := range vpcs {
		vpcIDs[i] = vpc.ID
	}

	// Scan subnets
	subnets, err := s.scanSubnets(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subnets: %w", err)
	}
	network.Subnets = subnets

	// Scan peering connections
	peeringConnections, err := s.scanPeeringConnections(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan peering connections: %w", err)
	}
	network.PeeringConnections = peeringConnections

	// Scan transit gateways
	transitGateways, err := s.scanTransitGateways(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to scan transit gateways: %w", err)
	}
	network.TransitGateways = transitGateways

	// Scan internet gateways
	internetGateways, err := s.scanInternetGateways(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan internet gateways: %w", err)
	}
	network.InternetGateways = internetGateways

	// Scan NAT gateways
	natGateways, err := s.scanNATGateways(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan NAT gateways: %w", err)
	}
	network.NATGateways = natGateways

	// Scan route tables
	routeTables, err := s.scanRouteTables(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan route tables: %w", err)
	}
	network.RouteTables = routeTables

	// Update subnet types based on route tables
	s.updateSubnetTypes(network)

	// Update VPC associations
	s.updateVPCAssociations(network)

	return network, nil
}

// scanVPCs scans VPCs
func (s *NetworkScanner) scanVPCs(ctx context.Context, vpcID string) ([]VPC, error) {
	input := &ec2.DescribeVpcsInput{}
	
	if vpcID != "" {
		input.VpcIds = []string{vpcID}
	}

	result, err := s.client.EC2.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, err
	}

	var vpcs []VPC
	for _, vpc := range result.Vpcs {
		v := VPC{
			ID:            *vpc.VpcId,
			CidrBlock:     *vpc.CidrBlock,
			State:         string(vpc.State),
			IsDefault:     vpc.IsDefault != nil && *vpc.IsDefault,
			DhcpOptionsID: *vpc.DhcpOptionsId,
			Tags:          convertTags(vpc.Tags),
		}
		
		// Get name from tags
		if name, ok := v.Tags["Name"]; ok {
			v.Name = name
		}
		
		vpcs = append(vpcs, v)
	}

	return vpcs, nil
}

// scanSubnets scans subnets
func (s *NetworkScanner) scanSubnets(ctx context.Context, vpcIDs []string) ([]Subnet, error) {
	if len(vpcIDs) == 0 {
		return []Subnet{}, nil
	}

	input := &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   &[]string{"vpc-id"}[0],
				Values: vpcIDs,
			},
		},
	}

	result, err := s.client.EC2.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, err
	}

	var subnets []Subnet
	for _, subnet := range result.Subnets {
		s := Subnet{
			ID:               *subnet.SubnetId,
			VpcID:            *subnet.VpcId,
			CidrBlock:        *subnet.CidrBlock,
			AvailabilityZone: *subnet.AvailabilityZone,
			State:            string(subnet.State),
			MapPublicIP:      subnet.MapPublicIpOnLaunch != nil && *subnet.MapPublicIpOnLaunch,
			Tags:             convertTags(subnet.Tags),
		}
		
		// Get name from tags
		if name, ok := s.Tags["Name"]; ok {
			s.Name = name
		}
		
		subnets = append(subnets, s)
	}

	return subnets, nil
}

// scanPeeringConnections scans VPC peering connections
func (s *NetworkScanner) scanPeeringConnections(ctx context.Context, vpcIDs []string) ([]PeeringConnection, error) {
	if len(vpcIDs) == 0 {
		return []PeeringConnection{}, nil
	}

	input := &ec2.DescribeVpcPeeringConnectionsInput{}

	result, err := s.client.EC2.DescribeVpcPeeringConnections(ctx, input)
	if err != nil {
		return nil, err
	}

	var connections []PeeringConnection
	for _, conn := range result.VpcPeeringConnections {
		// Only include connections involving our VPCs
		requesterVpcID := ""
		accepterVpcID := ""
		
		if conn.RequesterVpcInfo != nil && conn.RequesterVpcInfo.VpcId != nil {
			requesterVpcID = *conn.RequesterVpcInfo.VpcId
		}
		if conn.AccepterVpcInfo != nil && conn.AccepterVpcInfo.VpcId != nil {
			accepterVpcID = *conn.AccepterVpcInfo.VpcId
		}
		
		relevantConnection := false
		for _, vpcID := range vpcIDs {
			if vpcID == requesterVpcID || vpcID == accepterVpcID {
				relevantConnection = true
				break
			}
		}
		
		if !relevantConnection {
			continue
		}

		pc := PeeringConnection{
			ID:             *conn.VpcPeeringConnectionId,
			RequesterVpcID: requesterVpcID,
			AccepterVpcID:  accepterVpcID,
			Status:         string(conn.Status.Code),
			Tags:           convertTags(conn.Tags),
		}
		
		// Get name from tags
		if name, ok := pc.Tags["Name"]; ok {
			pc.Name = name
		}
		
		connections = append(connections, pc)
	}

	return connections, nil
}

// scanTransitGateways scans transit gateways
func (s *NetworkScanner) scanTransitGateways(ctx context.Context) ([]TransitGateway, error) {
	input := &ec2.DescribeTransitGatewaysInput{}

	result, err := s.client.EC2.DescribeTransitGateways(ctx, input)
	if err != nil {
		return nil, err
	}

	var tgws []TransitGateway
	for _, tgw := range result.TransitGateways {
		t := TransitGateway{
			ID:    *tgw.TransitGatewayId,
			State: string(tgw.State),
			Tags:  convertTags(tgw.Tags),
		}
		
		// Get name from tags
		if name, ok := t.Tags["Name"]; ok {
			t.Name = name
		}
		
		// Get attachments
		attachments, err := s.scanTransitGatewayAttachments(ctx, t.ID)
		if err != nil {
			// Log error but continue
			continue
		}
		t.Attachments = attachments
		
		tgws = append(tgws, t)
	}

	return tgws, nil
}

// scanTransitGatewayAttachments scans TGW attachments
func (s *NetworkScanner) scanTransitGatewayAttachments(ctx context.Context, tgwID string) ([]TransitGatewayAttachment, error) {
	input := &ec2.DescribeTransitGatewayAttachmentsInput{
		Filters: []types.Filter{
			{
				Name:   &[]string{"transit-gateway-id"}[0],
				Values: []string{tgwID},
			},
		},
	}

	result, err := s.client.EC2.DescribeTransitGatewayAttachments(ctx, input)
	if err != nil {
		return nil, err
	}

	var attachments []TransitGatewayAttachment
	for _, att := range result.TransitGatewayAttachments {
		a := TransitGatewayAttachment{
			ID:               *att.TransitGatewayAttachmentId,
			TransitGatewayID: *att.TransitGatewayId,
			ResourceType:     string(att.ResourceType),
			State:            string(att.State),
			Tags:             convertTags(att.Tags),
		}
		
		if att.ResourceId != nil {
			a.ResourceID = *att.ResourceId
		}
		
		attachments = append(attachments, a)
	}

	return attachments, nil
}

// scanInternetGateways scans internet gateways
func (s *NetworkScanner) scanInternetGateways(ctx context.Context, vpcIDs []string) ([]InternetGateway, error) {
	input := &ec2.DescribeInternetGatewaysInput{}

	result, err := s.client.EC2.DescribeInternetGateways(ctx, input)
	if err != nil {
		return nil, err
	}

	var igws []InternetGateway
	for _, igw := range result.InternetGateways {
		for _, attachment := range igw.Attachments {
			if attachment.VpcId == nil {
				continue
			}
			
			// Check if this IGW is attached to one of our VPCs
			vpcID := *attachment.VpcId
			relevantIGW := false
			for _, id := range vpcIDs {
				if id == vpcID {
					relevantIGW = true
					break
				}
			}
			
			if !relevantIGW {
				continue
			}
			
			ig := InternetGateway{
				ID:    *igw.InternetGatewayId,
				VpcID: vpcID,
				State: string(attachment.State),
				Tags:  convertTags(igw.Tags),
			}
			
			// Get name from tags
			if name, ok := ig.Tags["Name"]; ok {
				ig.Name = name
			}
			
			igws = append(igws, ig)
		}
	}

	return igws, nil
}

// scanNATGateways scans NAT gateways
func (s *NetworkScanner) scanNATGateways(ctx context.Context, vpcIDs []string) ([]NATGateway, error) {
	if len(vpcIDs) == 0 {
		return []NATGateway{}, nil
	}

	input := &ec2.DescribeNatGatewaysInput{}

	result, err := s.client.EC2.DescribeNatGateways(ctx, input)
	if err != nil {
		return nil, err
	}

	var natGws []NATGateway
	for _, nat := range result.NatGateways {
		// Filter by VPC ID
		if nat.VpcId == nil {
			continue
		}
		
		vpcID := *nat.VpcId
		relevantNAT := false
		for _, id := range vpcIDs {
			if id == vpcID {
				relevantNAT = true
				break
			}
		}
		
		if !relevantNAT {
			continue
		}
		
		ng := NATGateway{
			ID:               *nat.NatGatewayId,
			VpcID:            vpcID,
			SubnetID:         *nat.SubnetId,
			State:            string(nat.State),
			ConnectivityType: string(nat.ConnectivityType),
			Tags:             convertTags(nat.Tags),
		}
		
		// Get IP addresses
		for _, addr := range nat.NatGatewayAddresses {
			if addr.PublicIp != nil {
				ng.PublicIP = *addr.PublicIp
			}
			if addr.PrivateIp != nil {
				ng.PrivateIP = *addr.PrivateIp
			}
		}
		
		// Get name from tags
		if name, ok := ng.Tags["Name"]; ok {
			ng.Name = name
		}
		
		natGws = append(natGws, ng)
	}

	return natGws, nil
}

// scanRouteTables scans route tables
func (s *NetworkScanner) scanRouteTables(ctx context.Context, vpcIDs []string) ([]RouteTable, error) {
	if len(vpcIDs) == 0 {
		return []RouteTable{}, nil
	}

	input := &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{
				Name:   &[]string{"vpc-id"}[0],
				Values: vpcIDs,
			},
		},
	}

	result, err := s.client.EC2.DescribeRouteTables(ctx, input)
	if err != nil {
		return nil, err
	}

	var routeTables []RouteTable
	for _, rt := range result.RouteTables {
		r := RouteTable{
			ID:    *rt.RouteTableId,
			VpcID: *rt.VpcId,
			Tags:  convertTags(rt.Tags),
		}
		
		// Get name from tags
		if name, ok := r.Tags["Name"]; ok {
			r.Name = name
		}
		
		// Check if main route table
		for _, assoc := range rt.Associations {
			if assoc.Main != nil && *assoc.Main {
				r.IsMain = true
			}
			if assoc.SubnetId != nil {
				r.Associations = append(r.Associations, *assoc.SubnetId)
			}
		}
		
		// Get routes
		for _, route := range rt.Routes {
			ro := Route{
				State:  string(route.State),
				Origin: string(route.Origin),
			}
			
			if route.DestinationCidrBlock != nil {
				ro.DestinationCidr = *route.DestinationCidrBlock
			}
			if route.GatewayId != nil {
				ro.GatewayID = *route.GatewayId
			}
			if route.InstanceId != nil {
				ro.InstanceID = *route.InstanceId
			}
			if route.NetworkInterfaceId != nil {
				ro.NetworkInterfaceID = *route.NetworkInterfaceId
			}
			if route.VpcPeeringConnectionId != nil {
				ro.VpcPeeringID = *route.VpcPeeringConnectionId
			}
			if route.TransitGatewayId != nil {
				ro.TransitGatewayID = *route.TransitGatewayId
			}
			
			r.Routes = append(r.Routes, ro)
		}
		
		routeTables = append(routeTables, r)
	}

	return routeTables, nil
}

// updateSubnetTypes determines subnet types based on route tables
func (s *NetworkScanner) updateSubnetTypes(network *Network) {
	// Create a map of route table ID to route table
	routeTableMap := make(map[string]*RouteTable)
	for i := range network.RouteTables {
		routeTableMap[network.RouteTables[i].ID] = &network.RouteTables[i]
	}
	
	// Update each subnet
	for i := range network.Subnets {
		subnet := &network.Subnets[i]
		
		// Find route table for this subnet
		var routeTable *RouteTable
		for _, rt := range network.RouteTables {
			for _, assocSubnetID := range rt.Associations {
				if assocSubnetID == subnet.ID {
					routeTable = &rt
					subnet.RouteTableID = rt.ID
					break
				}
			}
			if routeTable != nil {
				break
			}
		}
		
		// If no explicit association, use main route table
		if routeTable == nil {
			for _, rt := range network.RouteTables {
				if rt.VpcID == subnet.VpcID && rt.IsMain {
					routeTable = &rt
					subnet.RouteTableID = rt.ID
					break
				}
			}
		}
		
		// Determine subnet type based on routes
		if routeTable != nil {
			subnet.Type = determineSubnetType(routeTable, network.InternetGateways)
		} else {
			subnet.Type = "isolated"
		}
	}
}

// determineSubnetType determines if a subnet is public, private, or isolated
func determineSubnetType(routeTable *RouteTable, igws []InternetGateway) string {
	hasIGWRoute := false
	hasNATRoute := false
	
	for _, route := range routeTable.Routes {
		// Check for internet gateway route
		if strings.HasPrefix(route.GatewayID, "igw-") {
			for _, igw := range igws {
				if igw.ID == route.GatewayID && route.DestinationCidr == "0.0.0.0/0" {
					hasIGWRoute = true
					break
				}
			}
		}
		
		// Check for NAT gateway route
		if strings.HasPrefix(route.GatewayID, "nat-") && route.DestinationCidr == "0.0.0.0/0" {
			hasNATRoute = true
		}
	}
	
	if hasIGWRoute {
		return "public"
	} else if hasNATRoute {
		return "private"
	}
	return "isolated"
}

// updateVPCAssociations updates VPC associations with subnets and other resources
func (s *NetworkScanner) updateVPCAssociations(network *Network) {
	// Create maps for quick lookup
	vpcMap := make(map[string]*VPC)
	for i := range network.VPCs {
		vpcMap[network.VPCs[i].ID] = &network.VPCs[i]
	}
	
	// Associate subnets with VPCs
	for _, subnet := range network.Subnets {
		if vpc, exists := vpcMap[subnet.VpcID]; exists {
			vpc.Subnets = append(vpc.Subnets, subnet.ID)
		}
	}
	
	// Associate internet gateways with VPCs
	for _, igw := range network.InternetGateways {
		if vpc, exists := vpcMap[igw.VpcID]; exists {
			vpc.InternetGateways = append(vpc.InternetGateways, igw.ID)
		}
	}
	
	// Associate NAT gateways with VPCs
	for _, nat := range network.NATGateways {
		if vpc, exists := vpcMap[nat.VpcID]; exists {
			vpc.NATGateways = append(vpc.NATGateways, nat.ID)
		}
	}
}

// convertTags converts AWS tags to map[string]string
func convertTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			result[*tag.Key] = *tag.Value
		}
	}
	return result
}
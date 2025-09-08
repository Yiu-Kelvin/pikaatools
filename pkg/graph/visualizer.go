package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

// Visualizer generates graph representations of AWS network infrastructure
type Visualizer struct {
	format string
}

// NewVisualizer creates a new graph visualizer
func NewVisualizer(format string) *Visualizer {
	return &Visualizer{
		format: format,
	}
}

// Generate generates a graph representation of the network
func (v *Visualizer) Generate(network *scanner.Network) (string, error) {
	switch v.format {
	case "text":
		return v.generateTextGraph(network), nil
	case "dot":
		return v.generateDotGraph(network), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", v.format)
	}
}

// generateTextGraph generates a text-based tree representation
func (v *Visualizer) generateTextGraph(network *scanner.Network) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("AWS Network Infrastructure - Region: %s\n", network.Region))
	result.WriteString(fmt.Sprintf("Scan Time: %s\n\n", network.ScanTime.Format("2006-01-02 15:04:05")))
	
	// Sort VPCs by ID for consistent output
	vpcs := make([]scanner.VPC, len(network.VPCs))
	copy(vpcs, network.VPCs)
	sort.Slice(vpcs, func(i, j int) bool {
		return vpcs[i].ID < vpcs[j].ID
	})
	
	// Create subnet map for quick lookup
	subnetMap := make(map[string]scanner.Subnet)
	for _, subnet := range network.Subnets {
		subnetMap[subnet.ID] = subnet
	}
	
	// Create peering map for quick lookup
	peeringMap := make(map[string][]scanner.PeeringConnection)
	for _, peering := range network.PeeringConnections {
		peeringMap[peering.RequesterVpcID] = append(peeringMap[peering.RequesterVpcID], peering)
		if peering.AccepterVpcID != peering.RequesterVpcID {
			peeringMap[peering.AccepterVpcID] = append(peeringMap[peering.AccepterVpcID], peering)
		}
	}
	
	// Create IGW map for quick lookup
	igwMap := make(map[string][]scanner.InternetGateway)
	for _, igw := range network.InternetGateways {
		igwMap[igw.VpcID] = append(igwMap[igw.VpcID], igw)
	}
	
	// Create NAT map for quick lookup
	natMap := make(map[string][]scanner.NATGateway)
	for _, nat := range network.NATGateways {
		natMap[nat.VpcID] = append(natMap[nat.VpcID], nat)
	}
	
	// Display VPCs and their resources
	for i, vpc := range vpcs {
		isLast := i == len(vpcs)-1
		v.writeVPC(&result, vpc, subnetMap, peeringMap, igwMap, natMap, isLast)
	}
	
	// Display Transit Gateways
	if len(network.TransitGateways) > 0 {
		result.WriteString("\n")
		for i, tgw := range network.TransitGateways {
			isLast := i == len(network.TransitGateways)-1
			v.writeTransitGateway(&result, tgw, network.VPCs, isLast)
		}
	}
	
	// Display summary
	result.WriteString(fmt.Sprintf("\nSummary:\n"))
	result.WriteString(fmt.Sprintf("  VPCs: %d\n", len(network.VPCs)))
	result.WriteString(fmt.Sprintf("  Subnets: %d\n", len(network.Subnets)))
	result.WriteString(fmt.Sprintf("  Peering Connections: %d\n", len(network.PeeringConnections)))
	result.WriteString(fmt.Sprintf("  Transit Gateways: %d\n", len(network.TransitGateways)))
	result.WriteString(fmt.Sprintf("  Internet Gateways: %d\n", len(network.InternetGateways)))
	result.WriteString(fmt.Sprintf("  NAT Gateways: %d\n", len(network.NATGateways)))
	
	return result.String()
}

// writeVPC writes a VPC and its associated resources
func (v *Visualizer) writeVPC(result *strings.Builder, vpc scanner.VPC, subnetMap map[string]scanner.Subnet, 
	peeringMap map[string][]scanner.PeeringConnection, igwMap map[string][]scanner.InternetGateway,
	natMap map[string][]scanner.NATGateway, isLastVPC bool) {
	
	vpcName := vpc.Name
	if vpcName == "" {
		vpcName = vpc.ID
	}
	
	defaultStr := ""
	if vpc.IsDefault {
		defaultStr = " [Default]"
	}
	
	result.WriteString(fmt.Sprintf("VPC: %s (%s)%s\n", vpcName, vpc.CidrBlock, defaultStr))
	
	// Count total items to display
	itemCount := 0
	itemCount += len(vpc.Subnets)
	if igws, exists := igwMap[vpc.ID]; exists {
		itemCount += len(igws)
	}
	if nats, exists := natMap[vpc.ID]; exists {
		itemCount += len(nats)
	}
	if peerings, exists := peeringMap[vpc.ID]; exists {
		itemCount += len(peerings)
	}
	
	currentItem := 0
	
	// Display subnets
	for _, subnetID := range vpc.Subnets {
		if subnet, exists := subnetMap[subnetID]; exists {
			currentItem++
			isLast := currentItem == itemCount
			v.writeSubnet(result, subnet, isLast)
		}
	}
	
	// Display Internet Gateways
	if igws, exists := igwMap[vpc.ID]; exists {
		for _, igw := range igws {
			currentItem++
			isLast := currentItem == itemCount
			v.writeInternetGateway(result, igw, isLast)
		}
	}
	
	// Display NAT Gateways
	if nats, exists := natMap[vpc.ID]; exists {
		for _, nat := range nats {
			currentItem++
			isLast := currentItem == itemCount
			v.writeNATGateway(result, nat, isLast)
		}
	}
	
	// Display Peering Connections
	if peerings, exists := peeringMap[vpc.ID]; exists {
		for _, peering := range peerings {
			currentItem++
			isLast := currentItem == itemCount
			v.writePeeringConnection(result, peering, vpc.ID, isLast)
		}
	}
	
	if !isLastVPC {
		result.WriteString("\n")
	}
}

// writeSubnet writes a subnet with proper tree formatting
func (v *Visualizer) writeSubnet(result *strings.Builder, subnet scanner.Subnet, isLast bool) {
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	
	subnetName := subnet.Name
	if subnetName == "" {
		subnetName = subnet.ID
	}
	
	typeStr := ""
	if subnet.Type != "" {
		typeStr = fmt.Sprintf(" [%s]", strings.Title(subnet.Type))
	}
	
	azStr := ""
	if subnet.AvailabilityZone != "" {
		azStr = fmt.Sprintf(" AZ:%s", subnet.AvailabilityZone)
	}
	
	result.WriteString(fmt.Sprintf("%sSubnet: %s (%s)%s%s\n", prefix, subnetName, subnet.CidrBlock, typeStr, azStr))
}

// writeInternetGateway writes an internet gateway
func (v *Visualizer) writeInternetGateway(result *strings.Builder, igw scanner.InternetGateway, isLast bool) {
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	
	igwName := igw.Name
	if igwName == "" {
		igwName = igw.ID
	}
	
	result.WriteString(fmt.Sprintf("%sInternet Gateway: %s [%s]\n", prefix, igwName, igw.State))
}

// writeNATGateway writes a NAT gateway
func (v *Visualizer) writeNATGateway(result *strings.Builder, nat scanner.NATGateway, isLast bool) {
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	
	natName := nat.Name
	if natName == "" {
		natName = nat.ID
	}
	
	ipInfo := ""
	if nat.PublicIP != "" {
		ipInfo = fmt.Sprintf(" Public:%s", nat.PublicIP)
	}
	if nat.PrivateIP != "" {
		ipInfo += fmt.Sprintf(" Private:%s", nat.PrivateIP)
	}
	
	result.WriteString(fmt.Sprintf("%sNAT Gateway: %s [%s]%s\n", prefix, natName, nat.State, ipInfo))
}

// writePeeringConnection writes a peering connection
func (v *Visualizer) writePeeringConnection(result *strings.Builder, peering scanner.PeeringConnection, currentVpcID string, isLast bool) {
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	
	peeringName := peering.Name
	if peeringName == "" {
		peeringName = peering.ID
	}
	
	// Determine the direction
	targetVPC := peering.AccepterVpcID
	direction := "→"
	if currentVpcID == peering.AccepterVpcID {
		targetVPC = peering.RequesterVpcID
		direction = "←"
	}
	
	result.WriteString(fmt.Sprintf("%sPeering: %s %s %s [%s]\n", prefix, peeringName, direction, targetVPC, peering.Status))
}

// writeTransitGateway writes a transit gateway and its attachments
func (v *Visualizer) writeTransitGateway(result *strings.Builder, tgw scanner.TransitGateway, vpcs []scanner.VPC, isLast bool) {
	tgwName := tgw.Name
	if tgwName == "" {
		tgwName = tgw.ID
	}
	
	result.WriteString(fmt.Sprintf("Transit Gateway: %s [%s]\n", tgwName, tgw.State))
	
	// Create VPC map for name lookup
	vpcMap := make(map[string]string)
	for _, vpc := range vpcs {
		name := vpc.Name
		if name == "" {
			name = vpc.ID
		}
		vpcMap[vpc.ID] = name
	}
	
	// Display attachments
	for i, attachment := range tgw.Attachments {
		isLastAttachment := i == len(tgw.Attachments)-1
		prefix := "├── "
		if isLastAttachment {
			prefix = "└── "
		}
		
		resourceName := attachment.ResourceID
		if attachment.ResourceType == "vpc" {
			if name, exists := vpcMap[attachment.ResourceID]; exists {
				resourceName = name
			}
		}
		
		result.WriteString(fmt.Sprintf("%sAttachment: %s (%s) [%s]\n", 
			prefix, resourceName, attachment.ResourceType, attachment.State))
	}
	
	if !isLast {
		result.WriteString("\n")
	}
}

// generateDotGraph generates a Graphviz DOT representation
func (v *Visualizer) generateDotGraph(network *scanner.Network) string {
	var result strings.Builder
	
	result.WriteString("digraph AWSNetwork {\n")
	result.WriteString("  rankdir=TB;\n")
	result.WriteString("  node [shape=box, style=rounded];\n")
	result.WriteString("  edge [fontsize=10];\n\n")
	
	// Define styles
	result.WriteString("  // Node styles\n")
	result.WriteString("  node [fillcolor=lightblue, style=\"rounded,filled\"];\n\n")
	
	// Add VPCs
	for _, vpc := range network.VPCs {
		vpcName := vpc.Name
		if vpcName == "" {
			vpcName = vpc.ID
		}
		
		label := fmt.Sprintf("%s\\n%s", vpcName, vpc.CidrBlock)
		if vpc.IsDefault {
			label += "\\n[Default]"
		}
		
		result.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=lightcyan];\n", vpc.ID, label))
	}
	
	// Add subnets
	result.WriteString("\n  // Subnets\n")
	for _, subnet := range network.Subnets {
		subnetName := subnet.Name
		if subnetName == "" {
			subnetName = subnet.ID
		}
		
		label := fmt.Sprintf("%s\\n%s\\n[%s]", subnetName, subnet.CidrBlock, strings.Title(subnet.Type))
		
		color := "lightgreen"
		switch subnet.Type {
		case "public":
			color = "lightgreen"
		case "private":
			color = "lightyellow"
		case "isolated":
			color = "lightcoral"
		}
		
		result.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=%s];\n", subnet.ID, label, color))
		result.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [style=dotted, label=\"contains\"];\n", subnet.VpcID, subnet.ID))
	}
	
	// Add Internet Gateways
	if len(network.InternetGateways) > 0 {
		result.WriteString("\n  // Internet Gateways\n")
		for _, igw := range network.InternetGateways {
			igwName := igw.Name
			if igwName == "" {
				igwName = igw.ID
			}
			
			result.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\nInternet Gateway\", fillcolor=orange];\n", igw.ID, igwName))
			result.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"attached\"];\n", igw.ID, igw.VpcID))
		}
	}
	
	// Add NAT Gateways
	if len(network.NATGateways) > 0 {
		result.WriteString("\n  // NAT Gateways\n")
		for _, nat := range network.NATGateways {
			natName := nat.Name
			if natName == "" {
				natName = nat.ID
			}
			
			label := fmt.Sprintf("%s\\nNAT Gateway", natName)
			if nat.PublicIP != "" {
				label += fmt.Sprintf("\\n%s", nat.PublicIP)
			}
			
			result.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=gold];\n", nat.ID, label))
			result.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [style=dotted, label=\"in\"];\n", nat.ID, nat.SubnetID))
		}
	}
	
	// Add peering connections
	if len(network.PeeringConnections) > 0 {
		result.WriteString("\n  // Peering Connections\n")
		for _, peering := range network.PeeringConnections {
			peeringName := peering.Name
			if peeringName == "" {
				peeringName = peering.ID
			}
			
			style := "solid"
			color := "blue"
			if peering.Status != "active" {
				style = "dashed"
				color = "gray"
			}
			
			result.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\\n[%s]\", style=%s, color=%s];\n", 
				peering.RequesterVpcID, peering.AccepterVpcID, peeringName, peering.Status, style, color))
		}
	}
	
	// Add Transit Gateways
	if len(network.TransitGateways) > 0 {
		result.WriteString("\n  // Transit Gateways\n")
		for _, tgw := range network.TransitGateways {
			tgwName := tgw.Name
			if tgwName == "" {
				tgwName = tgw.ID
			}
			
			result.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\nTransit Gateway\", fillcolor=purple, fontcolor=white];\n", tgw.ID, tgwName))
			
			// Add attachments
			for _, attachment := range tgw.Attachments {
				if attachment.ResourceType == "vpc" {
					style := "solid"
					if attachment.State != "available" {
						style = "dashed"
					}
					result.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"attached\", style=%s, color=purple];\n", 
						tgw.ID, attachment.ResourceID, style))
				}
			}
		}
	}
	
	result.WriteString("}\n")
	return result.String()
}
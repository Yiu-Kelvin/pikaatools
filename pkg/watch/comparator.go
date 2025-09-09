package watch

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/Yiu-Kelvin/pikaatools/pkg/scanner"
)

// Comparator compares two network states and reports differences
type Comparator struct {
	verbose bool
}

// NewComparator creates a new network state comparator
func NewComparator(verbose bool) *Comparator {
	return &Comparator{
		verbose: verbose,
	}
}

// LoadWorkingState loads a working state from a JSON file
func (c *Comparator) LoadWorkingState(filename string) (*scanner.Network, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read working state file %s: %w", filename, err)
	}

	var network scanner.Network
	err = json.Unmarshal(data, &network)
	if err != nil {
		return nil, fmt.Errorf("failed to parse working state JSON from %s: %w", filename, err)
	}

	return &network, nil
}

// Compare compares two network states and reports differences
func (c *Comparator) Compare(baseline, current *scanner.Network) []Difference {
	var differences []Difference

	// Compare VPCs
	differences = append(differences, c.compareVPCs(baseline.VPCs, current.VPCs)...)

	// Compare Subnets
	differences = append(differences, c.compareSubnets(baseline.Subnets, current.Subnets)...)

	// Compare Security Groups
	differences = append(differences, c.compareSecurityGroups(baseline.SecurityGroups, current.SecurityGroups)...)

	// Compare Network ACLs
	differences = append(differences, c.compareNetworkAcls(baseline.NetworkAcls, current.NetworkAcls)...)

	// Compare Route Tables
	differences = append(differences, c.compareRouteTables(baseline.RouteTables, current.RouteTables)...)

	// Compare Peering Connections
	differences = append(differences, c.comparePeeringConnections(baseline.PeeringConnections, current.PeeringConnections)...)

	// Compare Transit Gateways
	differences = append(differences, c.compareTransitGateways(baseline.TransitGateways, current.TransitGateways)...)

	// Compare Internet Gateways
	differences = append(differences, c.compareInternetGateways(baseline.InternetGateways, current.InternetGateways)...)

	// Compare NAT Gateways
	differences = append(differences, c.compareNATGateways(baseline.NATGateways, current.NATGateways)...)

	// Compare IAM Roles
	differences = append(differences, c.compareIAMRoles(baseline.IAMRoles, current.IAMRoles)...)

	return differences
}

// PrintDifferences prints differences in colored output
func (c *Comparator) PrintDifferences(differences []Difference) {
	if len(differences) == 0 {
		color.Green("✓ No differences found - infrastructure state matches baseline")
		return
	}

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s %s\n", red("⚠"), red(fmt.Sprintf("Found %d differences:", len(differences))))
	fmt.Println()

	for _, diff := range differences {
		switch diff.Type {
		case Added:
			fmt.Printf("%s %s: %s %s\n", red("+ ADDED"), cyan(diff.ResourceType), yellow(diff.ResourceID), diff.Description)
		case Removed:
			fmt.Printf("%s %s: %s %s\n", red("- REMOVED"), cyan(diff.ResourceType), yellow(diff.ResourceID), diff.Description)
		case Modified:
			fmt.Printf("%s %s: %s %s\n", red("~ MODIFIED"), cyan(diff.ResourceType), yellow(diff.ResourceID), diff.Description)
		}

		if c.verbose && len(diff.Details) > 0 {
			for _, detail := range diff.Details {
				fmt.Printf("    %s\n", detail)
			}
		}
	}
	fmt.Println()
}

// Difference represents a difference between two network states
type Difference struct {
	Type         DifferenceType
	ResourceType string
	ResourceID   string
	Description  string
	Details      []string
}

// DifferenceType represents the type of difference
type DifferenceType int

const (
	Added DifferenceType = iota
	Removed
	Modified
)

// Helper functions for comparing different resource types
func (c *Comparator) compareVPCs(baseline, current []scanner.VPC) []Difference {
	return c.compareSlices("VPC", baseline, current, func(v interface{}) string { 
		return v.(scanner.VPC).ID 
	})
}

func (c *Comparator) compareSubnets(baseline, current []scanner.Subnet) []Difference {
	return c.compareSlices("Subnet", baseline, current, func(s interface{}) string { 
		return s.(scanner.Subnet).ID 
	})
}

func (c *Comparator) compareSecurityGroups(baseline, current []scanner.SecurityGroup) []Difference {
	return c.compareSlices("SecurityGroup", baseline, current, func(sg interface{}) string { 
		return sg.(scanner.SecurityGroup).ID 
	})
}

func (c *Comparator) compareNetworkAcls(baseline, current []scanner.NetworkAcl) []Difference {
	return c.compareSlices("NetworkACL", baseline, current, func(nacl interface{}) string { 
		return nacl.(scanner.NetworkAcl).ID 
	})
}

func (c *Comparator) compareRouteTables(baseline, current []scanner.RouteTable) []Difference {
	return c.compareSlices("RouteTable", baseline, current, func(rt interface{}) string { 
		return rt.(scanner.RouteTable).ID 
	})
}

func (c *Comparator) comparePeeringConnections(baseline, current []scanner.PeeringConnection) []Difference {
	return c.compareSlices("PeeringConnection", baseline, current, func(pc interface{}) string { 
		return pc.(scanner.PeeringConnection).ID 
	})
}

func (c *Comparator) compareTransitGateways(baseline, current []scanner.TransitGateway) []Difference {
	return c.compareSlices("TransitGateway", baseline, current, func(tgw interface{}) string { 
		return tgw.(scanner.TransitGateway).ID 
	})
}

func (c *Comparator) compareInternetGateways(baseline, current []scanner.InternetGateway) []Difference {
	return c.compareSlices("InternetGateway", baseline, current, func(igw interface{}) string { 
		return igw.(scanner.InternetGateway).ID 
	})
}

func (c *Comparator) compareNATGateways(baseline, current []scanner.NATGateway) []Difference {
	return c.compareSlices("NATGateway", baseline, current, func(nat interface{}) string { 
		return nat.(scanner.NATGateway).ID 
	})
}

func (c *Comparator) compareIAMRoles(baseline, current []scanner.IAMRole) []Difference {
	return c.compareSlices("IAMRole", baseline, current, func(role interface{}) string { 
		return role.(scanner.IAMRole).ID 
	})
}

// Generic slice comparison function  
func (c *Comparator) compareSlices(resourceType string, baseline, current interface{}, getID func(interface{}) string) []Difference {
	var differences []Difference

	// Use reflection to handle the interface{} types
	baselineSlice := reflect.ValueOf(baseline)
	currentSlice := reflect.ValueOf(current)

	if baselineSlice.Kind() != reflect.Slice || currentSlice.Kind() != reflect.Slice {
		return differences
	}

	// Create maps for quick lookup
	baselineMap := make(map[string]interface{})
	currentMap := make(map[string]interface{})

	for i := 0; i < baselineSlice.Len(); i++ {
		item := baselineSlice.Index(i).Interface()
		id := getID(item)
		baselineMap[id] = item
	}

	for i := 0; i < currentSlice.Len(); i++ {
		item := currentSlice.Index(i).Interface()
		id := getID(item)
		currentMap[id] = item
	}

	// Find added items
	for id := range currentMap {
		if _, exists := baselineMap[id]; !exists {
			differences = append(differences, Difference{
				Type:         Added,
				ResourceType: resourceType,
				ResourceID:   id,
				Description:  fmt.Sprintf("New %s created", strings.ToLower(resourceType)),
			})
		}
	}

	// Find removed items
	for id := range baselineMap {
		if _, exists := currentMap[id]; !exists {
			differences = append(differences, Difference{
				Type:         Removed,
				ResourceType: resourceType,
				ResourceID:   id,
				Description:  fmt.Sprintf("%s was deleted", strings.ToLower(resourceType)),
			})
		}
	}

	// Find modified items
	for id, currentItem := range currentMap {
		if baselineItem, exists := baselineMap[id]; exists {
			if details := c.findObjectDifferences(baselineItem, currentItem); len(details) > 0 {
				differences = append(differences, Difference{
					Type:         Modified,
					ResourceType: resourceType,
					ResourceID:   id,
					Description:  fmt.Sprintf("%s configuration changed", strings.ToLower(resourceType)),
					Details:      details,
				})
			}
		}
	}

	return differences
}

// findObjectDifferences compares two objects and returns a list of field differences
func (c *Comparator) findObjectDifferences(baseline, current interface{}) []string {
	var details []string

	baselineValue := reflect.ValueOf(baseline)
	currentValue := reflect.ValueOf(current)

	if baselineValue.Type() != currentValue.Type() {
		return []string{"Objects have different types"}
	}

	switch baselineValue.Kind() {
	case reflect.Struct:
		details = append(details, c.compareStructs(baselineValue, currentValue, "")...)
	case reflect.Slice:
		details = append(details, c.compareSlicesReflect(baselineValue, currentValue, "")...)
	case reflect.Map:
		details = append(details, c.compareMaps(baselineValue, currentValue, "")...)
	default:
		if !reflect.DeepEqual(baseline, current) {
			details = append(details, fmt.Sprintf("Value changed from %v to %v", baseline, current))
		}
	}

	return details
}

func (c *Comparator) compareStructs(baseline, current reflect.Value, path string) []string {
	var details []string
	structType := baseline.Type()

	for i := 0; i < baseline.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name

		// Skip private fields and certain fields we don't want to compare
		if !field.IsExported() || c.shouldSkipField(fieldName) {
			continue
		}

		fieldPath := fieldName
		if path != "" {
			fieldPath = fmt.Sprintf("%s.%s", path, fieldName)
		}

		baselineField := baseline.Field(i)
		currentField := current.Field(i)

		if !reflect.DeepEqual(baselineField.Interface(), currentField.Interface()) {
			switch baselineField.Kind() {
			case reflect.Struct:
				details = append(details, c.compareStructs(baselineField, currentField, fieldPath)...)
			case reflect.Slice:
				details = append(details, c.compareSlicesReflect(baselineField, currentField, fieldPath)...)
			case reflect.Map:
				details = append(details, c.compareMaps(baselineField, currentField, fieldPath)...)
			default:
				details = append(details, fmt.Sprintf("%s: %v → %v", fieldPath, baselineField.Interface(), currentField.Interface()))
			}
		}
	}

	return details
}

func (c *Comparator) compareSlicesReflect(baseline, current reflect.Value, path string) []string {
	var details []string

	if baseline.Len() != current.Len() {
		details = append(details, fmt.Sprintf("%s: length changed from %d to %d", path, baseline.Len(), current.Len()))
	}

	// For simplicity, we'll just note that the slice changed if lengths differ
	// More sophisticated comparison could be implemented for specific slice types
	if baseline.Len() != current.Len() || !reflect.DeepEqual(baseline.Interface(), current.Interface()) {
		details = append(details, fmt.Sprintf("%s: slice contents changed", path))
	}

	return details
}

func (c *Comparator) compareMaps(baseline, current reflect.Value, path string) []string {
	var details []string

	// Check for added/removed keys
	for _, key := range baseline.MapKeys() {
		if !current.MapIndex(key).IsValid() {
			details = append(details, fmt.Sprintf("%s[%v]: key removed", path, key.Interface()))
		}
	}

	for _, key := range current.MapKeys() {
		baselineValue := baseline.MapIndex(key)
		currentValue := current.MapIndex(key)

		if !baselineValue.IsValid() {
			details = append(details, fmt.Sprintf("%s[%v]: key added with value %v", path, key.Interface(), currentValue.Interface()))
		} else if !reflect.DeepEqual(baselineValue.Interface(), currentValue.Interface()) {
			details = append(details, fmt.Sprintf("%s[%v]: %v → %v", path, key.Interface(), baselineValue.Interface(), currentValue.Interface()))
		}
	}

	return details
}

// shouldSkipField determines if a field should be skipped during comparison
func (c *Comparator) shouldSkipField(fieldName string) bool {
	skipFields := []string{"ScanTime", "CreateDate", "UpdateDate"}
	for _, skip := range skipFields {
		if fieldName == skip {
			return true
		}
	}
	return false
}
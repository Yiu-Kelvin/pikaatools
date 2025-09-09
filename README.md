# PikaaTools - AWS Network Scanner

A Go-based tool that scans your AWS network infrastructure and visualizes it as a graph. This tool discovers VPCs, subnets, peering connections, Transit Gateways, and other network resources to help you understand your AWS network topology.

## Features

- 🔍 **Comprehensive Scanning**: Discovers VPCs, subnets, peering connections, Transit Gateways, route tables, security groups with detailed rules, Network ACLs with entries, IAM roles and policies, and more
- 👀 **Change Watching**: Monitor infrastructure changes with `watch` command that compares current state against a baseline and highlights differences in red
- 📊 **Graph Visualization**: Generates text-based network topology graphs
- 💾 **JSON Export**: Save complete working state to JSON file for analysis and automation
- 🔧 **Configurable**: Support for multiple AWS profiles and regions
- 🚀 **Fast**: Concurrent scanning for efficient discovery
- 🔒 **Secure**: Uses standard AWS credential chain
- 📝 **Verbose Mode**: Detailed timing information for each resource scan

## Installation

```bash
go install github.com/Yiu-Kelvin/pikaatools@latest
```

Or build from source:

```bash
git clone https://github.com/Yiu-Kelvin/pikaatools.git
cd pikaatools
go build -o pikaatools .
```

## Usage

### Basic Usage

```bash
# Scan all VPCs in default region
./pikaatools scan

# Scan specific region
./pikaatools scan --region us-west-2

# Scan with specific AWS profile
./pikaatools scan --profile myprofile

# Output graph in DOT format
./pikaatools scan --output dot

# Scan specific VPC
./pikaatools scan --vpc-id vpc-12345678

# Export working state to JSON file
./pikaatools scan --export-json my_network.json

# Save working state to default file (working_state.json)
./pikaatools scan --save-state

# Enable verbose output with timing information
./pikaatools scan --verbose

# Combine flags for detailed verbose scanning of specific VPC
./pikaatools scan --vpc-id vpc-12345678 --verbose --export-json detailed_scan.json
```

### Watch for Changes

```bash
# First, create a baseline working state
./pikaatools scan --save-state

# Watch for changes against the baseline (scans every 30 seconds by default)
./pikaatools watch

# Watch with custom interval and specific working state file
./pikaatools watch -f my_baseline.json --interval 1m

# Watch with verbose output for detailed scanning information
./pikaatools watch --verbose --interval 2m

# Watch specific VPC for changes
./pikaatools watch --vpc-id vpc-12345678 --interval 45s
```

### Configuration

The tool uses the standard AWS credential chain:
1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. AWS credentials file (`~/.aws/credentials`)
3. IAM roles for EC2 instances
4. IAM roles for ECS tasks

### Required Permissions

The tool requires the following AWS permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeVpcs",
                "ec2:DescribeSubnets",
                "ec2:DescribeVpcPeeringConnections",
                "ec2:DescribeTransitGateways",
                "ec2:DescribeTransitGatewayAttachments",
                "ec2:DescribeRouteTables",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeNatGateways",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeNetworkAcls",
                "ec2:DescribeNetworkAcls",
                "iam:ListRoles",
                "iam:GetRole",
                "iam:ListAttachedRolePolicies",
                "iam:ListRolePolicies",
                "iam:GetRolePolicy",
                "iam:GetPolicy",
                "iam:GetPolicyVersion"
            ],
            "Resource": "*"
        }
    ]
}
```

## Output Formats

### Text Graph (Default)
```
VPC: vpc-12345678 (10.0.0.0/16)
├── Subnet: subnet-abc123 (10.0.1.0/24) [Public]
├── Subnet: subnet-def456 (10.0.2.0/24) [Private]
└── Peering: pcx-789xyz → vpc-87654321

Transit Gateway: tgw-12345678
├── Attachment: vpc-12345678
└── Attachment: vpc-87654321
```

### DOT Format
Generate Graphviz DOT files for advanced visualization:

```bash
./pikaatools scan --output dot > network.dot
dot -Tpng network.dot -o network.png
```

### JSON Format
Export complete network state for analysis, automation, or integration:

```bash
./pikaatools scan --save-state
```

This creates a `working_state.json` file containing all discovered resources with their complete configurations including:
- VPCs with CIDR blocks, tags, and associated resources
- Subnets with availability zones, route tables, Network ACL associations, and types (public/private/isolated)
- Security groups with detailed inbound and outbound rules, including protocols, ports, CIDR blocks, and referenced security groups
- Network ACLs with entries including rule numbers, protocols, actions, port ranges, and ICMP types
- Route tables with all routes and associations
- Transit Gateways with attachments
- Internet Gateways and NAT Gateways
- VPC Peering connections
- IAM roles with attached and inline policies

### Verbose Mode

Enable verbose output to see detailed timing information for each resource scan:

```bash
./pikaatools scan --verbose
```

Example verbose output:
```
Initializing AWS client...
Scanning AWS network infrastructure in region: us-east-1
Scanned vpc vpc-12345678 took 15.2ms
Scanned vpc vpc-87654321 took 12.8ms
Scanned 2 VPCs took 28ms
Scanned 4 subnets took 45ms
Scanned 5 security groups took 78ms
Scanned 3 network ACLs took 42ms
Scanned 12 IAM roles took 156ms
Found 2 VPCs, 4 subnets, 1 peering connections, 0 transit gateways, 5 security groups, 3 network ACLs, 12 IAM roles
```

### Watch Mode

Monitor your infrastructure for changes in real-time:

```bash
./pikaatools watch --verbose
```

Example watch output:
```
Loading baseline state from working_state.json...
Loaded baseline state from working_state.json (scanned at 2024-01-15T10:30:00Z)
Starting periodic scan every 30s...

🔍 Starting initial scan...
[2024-01-15 10:35:00] Scan completed in 2.3s (region: us-east-1)
✓ No differences found - infrastructure state matches baseline

🔍 Performing periodic scan...
[2024-01-15 10:35:30] Scan completed in 1.8s (region: us-east-1)
⚠ Found 2 differences:

+ ADDED VPC: vpc-new123 New vpc created
~ MODIFIED SecurityGroup: sg-12345 security group configuration changed
    IngressRules: slice contents changed
    EgressRules: slice contents changed
```
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- AWS SDK for Go team
- Graphviz for visualization capabilities
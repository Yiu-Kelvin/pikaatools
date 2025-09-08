# PikaaTools - AWS Network Scanner

A Go-based tool that scans your AWS network infrastructure and visualizes it as a graph. This tool discovers VPCs, subnets, peering connections, Transit Gateways, and other network resources to help you understand your AWS network topology.

## Features

- 🔍 **Comprehensive Scanning**: Discovers VPCs, subnets, peering connections, Transit Gateways, route tables, security groups, and more
- 📊 **Graph Visualization**: Generates text-based network topology graphs
- 💾 **JSON Export**: Save complete working state to JSON file for analysis and automation
- 🔧 **Configurable**: Support for multiple AWS profiles and regions
- 🚀 **Fast**: Concurrent scanning for efficient discovery
- 🔒 **Secure**: Uses standard AWS credential chain

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
                "ec2:DescribeNetworkAcls"
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
- Subnets with availability zones, route tables, and types (public/private/isolated)
- Security groups with rules and associations
- Route tables with all routes and associations
- Transit Gateways with attachments
- Internet Gateways and NAT Gateways
- VPC Peering connections
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
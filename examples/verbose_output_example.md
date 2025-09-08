# Verbose Output Example

This example shows the enhanced verbose output with timing information and security group rules.

## Basic Usage

```bash
./pikaatools scan --verbose --region us-east-1
```

## Expected Verbose Output

```
Initializing AWS client...
Scanning AWS network infrastructure in region: us-east-1
Scanned vpc vpc-12345678 took 15.2ms
Scanned vpc vpc-87654321 took 12.8ms
Scanned 2 VPCs took 28ms
Scanned 4 subnets took 45ms
Scanned 1 peering connections took 22ms
Scanned 0 transit gateways took 18ms
Scanned 2 internet gateways took 35ms
Scanned 1 NAT gateways took 28ms
Scanned 3 route tables took 52ms
Scanned 5 security groups took 78ms
Scanned 12 IAM roles took 156ms
Found 2 VPCs, 4 subnets, 1 peering connections, 0 transit gateways, 5 security groups, 12 IAM roles

VPC: vpc-12345678 (10.0.0.0/16)
├── Subnet: subnet-abc123 (10.0.1.0/24) [Public]
├── Subnet: subnet-def456 (10.0.2.0/24) [Private]
└── Peering: pcx-789xyz → vpc-87654321
```

## Security Group Rules Enhancement

The security groups now include detailed rule information in the JSON export:

```json
{
  "security_groups": [
    {
      "id": "sg-12345678",
      "name": "web-server-sg",
      "description": "Security group for web servers",
      "vpc_id": "vpc-12345678",
      "tags": {
        "Name": "web-server-sg",
        "Environment": "production"
      },
      "ingress_rules": [
        {
          "ip_protocol": "tcp",
          "from_port": 80,
          "to_port": 80,
          "cidr_blocks": ["0.0.0.0/0"],
          "ipv6_cidr_blocks": [],
          "prefix_list_ids": [],
          "referenced_group_id": "",
          "referenced_group_owner_id": "",
          "description": "Allow HTTP traffic",
          "tags": {}
        },
        {
          "ip_protocol": "tcp",
          "from_port": 443,
          "to_port": 443,
          "cidr_blocks": ["0.0.0.0/0"],
          "ipv6_cidr_blocks": [],
          "prefix_list_ids": [],
          "referenced_group_id": "",
          "referenced_group_owner_id": "",
          "description": "Allow HTTPS traffic",
          "tags": {}
        }
      ],
      "egress_rules": [
        {
          "ip_protocol": "-1",
          "from_port": -1,
          "to_port": -1,
          "cidr_blocks": ["0.0.0.0/0"],
          "ipv6_cidr_blocks": [],
          "prefix_list_ids": [],
          "referenced_group_id": "",
          "referenced_group_owner_id": "",
          "description": "Allow all outbound traffic",
          "tags": {}
        }
      ]
    }
  ]
}
```

## Features Added

1. **Timing Information**: Each resource type scan shows timing information
2. **Individual Resource Timing**: VPCs show individual scan times (e.g., "Scanned vpc vpc-12345678 took 15.2ms")
3. **Security Group Rules**: Complete ingress and egress rule details including:
   - IP protocol (tcp, udp, icmp, etc.)
   - Port ranges (from_port, to_port)
   - CIDR blocks (IPv4 and IPv6)
   - Prefix list IDs
   - Referenced security group IDs
   - Rule descriptions
   - Rule tags
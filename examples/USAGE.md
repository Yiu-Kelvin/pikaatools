# Example configurations and usage

## Environment Variables

```bash
# AWS Credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1

# Optional: AWS Session Token for temporary credentials
export AWS_SESSION_TOKEN=your_session_token
```

## AWS Credentials File

Create `~/.aws/credentials`:

```ini
[default]
aws_access_key_id = your_access_key
aws_secret_access_key = your_secret_key

[production]
aws_access_key_id = prod_access_key
aws_secret_access_key = prod_secret_key

[development]
aws_access_key_id = dev_access_key
aws_secret_access_key = dev_secret_key
```

Create `~/.aws/config`:

```ini
[default]
region = us-east-1
output = json

[profile production]
region = us-west-2
output = json

[profile development]
region = us-east-1
output = json
```

## Example Usage Scenarios

### Basic Scanning

```bash
# Scan all VPCs in default region with default profile
./pikaatools scan

# Scan with verbose output
./pikaatools scan --verbose

# Scan specific region
./pikaatools scan --region us-west-2

# Scan using specific AWS profile
./pikaatools scan --profile production
```

### Targeted Scanning

```bash
# Scan specific VPC
./pikaatools scan --vpc-id vpc-12345678

# Scan with multiple filters
./pikaatools scan --region us-west-2 --profile production --verbose
```

### Output Formats

```bash
# Default text output
./pikaatools scan > network-topology.txt

# Generate DOT file for Graphviz
./pikaatools scan --output dot > network.dot

# Convert DOT to PNG (requires Graphviz installed)
./pikaatools scan --output dot | dot -Tpng -o network.png

# Convert DOT to SVG
./pikaatools scan --output dot | dot -Tsvg -o network.svg

# Convert DOT to PDF
./pikaatools scan --output dot | dot -Tpdf -o network.pdf
```

## Sample Output

### Text Output
```
AWS Network Infrastructure - Region: us-east-1
Scan Time: 2025-01-15 10:30:45

VPC: Production VPC (10.0.0.0/16)
├── Subnet: Public Subnet 1 (10.0.1.0/24) [Public] AZ:us-east-1a
├── Subnet: Private Subnet 1 (10.0.10.0/24) [Private] AZ:us-east-1a
├── Subnet: Private Subnet 2 (10.0.11.0/24) [Private] AZ:us-east-1b
├── Internet Gateway: prod-igw [available]
├── NAT Gateway: prod-nat-1 [available] Public:1.2.3.4 Private:10.0.1.5
└── Peering: prod-to-dev → vpc-87654321 [active]

VPC: Development VPC (172.16.0.0/16)
└── Subnet: Dev Subnet (172.16.1.0/24) [Public] AZ:us-east-1a

Transit Gateway: Main TGW [available]
├── Attachment: Production VPC (vpc) [available]
└── Attachment: Development VPC (vpc) [available]

Summary:
  VPCs: 2
  Subnets: 4
  Peering Connections: 1
  Transit Gateways: 1
  Internet Gateways: 1
  NAT Gateways: 1
```

## Graphviz Installation

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install graphviz
```

### macOS
```bash
brew install graphviz
```

### Windows
Download from: https://graphviz.org/download/

## Integration with Other Tools

### CI/CD Pipeline

```yaml
# .github/workflows/network-scan.yml
name: Network Documentation
on:
  schedule:
    - cron: '0 6 * * 1'  # Run weekly on Monday
  workflow_dispatch:

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Install Graphviz
        run: sudo apt-get install graphviz
      - name: Build scanner
        run: go build -o pikaatools .
      - name: Scan network
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: us-east-1
        run: |
          ./pikaatools scan --output dot > docs/network.dot
          dot -Tpng docs/network.dot -o docs/network.png
          dot -Tsvg docs/network.dot -o docs/network.svg
      - name: Commit documentation
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add docs/
          git commit -m "Update network documentation" || exit 0
          git push
```

### Monitoring Integration

```bash
#!/bin/bash
# network-check.sh

# Scan network and send to monitoring system
SCAN_RESULT=$(./pikaatools scan --region us-east-1 2>&1)
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    # Send alert to monitoring system
    curl -X POST "https://hooks.slack.com/webhook-url" \
         -H 'Content-type: application/json' \
         --data "{\"text\":\"Network scan failed: $SCAN_RESULT\"}"
else
    echo "Network scan completed successfully"
fi
```
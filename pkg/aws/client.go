package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// Client wraps AWS services needed for network scanning
type Client struct {
	EC2    *ec2.Client
	IAM    *iam.Client
	config aws.Config
}

// NewClient creates a new AWS client with the specified region and profile
func NewClient(ctx context.Context, region, profile string) (*Client, error) {
	var opts []func(*config.LoadOptions) error
	
	// Set region
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1" // Default region
		}
	}
	opts = append(opts, config.WithRegion(region))
	
	// Set profile if specified
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}
	
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	
	return &Client{
		EC2:    ec2.NewFromConfig(cfg),
		IAM:    iam.NewFromConfig(cfg),
		config: cfg,
	}, nil
}

// Region returns the current AWS region
func (c *Client) Region() string {
	return c.config.Region
}
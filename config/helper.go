package config

// helper functions
import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/natemarks/vpc_flow_logs/version"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// GetLogger returns a logger for the application
func GetLogger(debug bool) (log zerolog.Logger) {
	log = zerolog.New(os.Stdout).With().Str("version", version.Version).Timestamp().Logger()
	log = log.Level(zerolog.InfoLevel)
	if debug {
		log = log.Level(zerolog.DebugLevel)
		log.Debug().Msg("debug logging enabled")
	}
	return log
}

// AWSInfo represents information about the AWS environment
type AWSInfo struct {
	Region    string `json:"region"`
	AccountID string `json:"accountID"`
}

// GetAWSInfo retrieves the current AWS region and account ID
func GetAWSInfo() (*AWSInfo, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create STS client
	client := sts.NewFromConfig(cfg)

	// Retrieve caller identity to get account ID
	identityOutput, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving caller identity: %v", err)
	}

	// Construct AWSInfo object
	awsInfo := &AWSInfo{
		Region:    cfg.Region,
		AccountID: *identityOutput.Account,
	}

	return awsInfo, nil
}

// RoleNameFromVPC returns the name of the IAM role for a VPC
func RoleNameFromVPC(vpcID string) string {
	return fmt.Sprintf("%v_flow_log_to_cloudwatch_logs", vpcID)

}

// PolicyARNFromName returns the ARN of an IAM policy given its name
func PolicyARNFromName(policyName string, awsInfo AWSInfo) string {
	return fmt.Sprintf("arn:aws:iam::%v:policy/%v", awsInfo.AccountID, policyName)
}

// DetachRolePolicy detaches an IAM policy from a role
func DetachRolePolicy(roleARN, policyARN string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Detach policy
	_, err = client.DetachRolePolicy(context.TODO(), &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyARN),
		RoleName:  aws.String(roleARN),
	})
	if err != nil {
		return fmt.Errorf("error detaching policy from role: %v", err)
	}

	return nil
}

// DeleteRole deletes an IAM role
func DeleteRole(roleName string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Delete role
	_, err = client.DeleteRole(context.TODO(), &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}

	return nil
}

// DeletePolicy deletes an IAM policy
func DeletePolicy(policyARN string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Delete policy input
	deletePolicyInput := &iam.DeletePolicyInput{
		PolicyArn: &policyARN,
	}

	_, err = client.DeletePolicy(context.TODO(), deletePolicyInput)
	if err != nil {
		return fmt.Errorf("error deleting policy: %v", err)
	}

	return nil
}

// DeleteCloudwatchLogGroup deletes a CloudWatch log group
func DeleteCloudwatchLogGroup(logGroupName string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create CloudWatch Logs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	// Delete log group
	_, err = client.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		return fmt.Errorf("error deleting log group: %v", err)
	}

	return nil
}

// DeleteVPCFlowLog delete a VPC flow log
func DeleteVPCFlowLog(flowLogID string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create EC2 client
	client := ec2.NewFromConfig(cfg)

	// Delete flow log input
	deleteFlowLogInput := &ec2.DeleteFlowLogsInput{
		FlowLogIds: []string{flowLogID},
	}

	_, err = client.DeleteFlowLogs(context.TODO(), deleteFlowLogInput)
	if err != nil {
		return fmt.Errorf("error deleting flow log: %v", err)
	}

	return nil
}

// GetLogGroupARN retrieves the ARN of a CloudWatch log group given its name
func GetLogGroupARN(logGroupName string) (string, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	// Create CloudWatch Logs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	// Describe log group
	describeLogGroupsOutput, err := client.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if err != nil {
		return "", fmt.Errorf("error describing log groups: %v", err)
	}

	// Check if the log group with the given name exists
	if len(describeLogGroupsOutput.LogGroups) == 0 {
		return "", fmt.Errorf("log group '%s' not found", logGroupName)
	}

	// Return the ARN of the first log group found (assuming there's only one)
	return *describeLogGroupsOutput.LogGroups[0].Arn, nil
}

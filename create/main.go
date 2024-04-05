package create

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	appcfg "github.com/natemarks/vpc_flow_logs/config"
	apptypes "github.com/natemarks/vpc_flow_logs/types"
)

const (
	defaultRetentionDays int32 = 14
)

func logGroupNameFromVPC(vpcID string) string {
	return fmt.Sprintf("/aws/vpc/%s/flowlogs", vpcID)

}

func logGroupExists(mylg apptypes.LogGroup) bool {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("error loading AWS SDK configuration")
	}

	// Create CloudWatch Logs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	// Describe log group input
	describeLogGroupsInput := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(mylg.Name),
	}

	// Describe log groups
	describeLogGroupsOutput, err := client.DescribeLogGroups(context.TODO(), describeLogGroupsInput)
	if err != nil {
		panic("error describing CloudWatch log groups")
	}

	// Check if the log group with the given name exists
	for _, lg := range describeLogGroupsOutput.LogGroups {
		if *lg.LogGroupName == mylg.Name {
			return true
		}
	}

	// Log group not found
	return false
}

func logGroup(cConfig appcfg.CreateConfig) (lg *apptypes.LogGroup, err error) {
	var result apptypes.LogGroup
	result.Name = logGroupNameFromVPC(cConfig.VPCID)
	result.ARN = fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", cConfig.Region, cConfig.AccountID, result.Name)
	result.RetentionDays = int(defaultRetentionDays)

	if logGroupExists(result) {
		return &result, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return lg, fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}

	client := cloudwatchlogs.NewFromConfig(cfg)

	createLogGroupInput := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(result.Name),
	}

	_, err = client.CreateLogGroup(context.TODO(), createLogGroupInput)
	if err != nil {
		return nil, fmt.Errorf("error creating CloudWatch log group: %v", err)
	}

	putRetentionPolicyInput := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(result.Name),
		RetentionInDays: aws.Int32(defaultRetentionDays),
	}

	_, err = client.PutRetentionPolicy(context.TODO(), putRetentionPolicyInput)
	if err != nil {
		return nil, fmt.Errorf("error setting retention policy for CloudWatch log group: %v", err)
	}
	lg = &result
	return lg, nil
}

func policyDocument(lg apptypes.LogGroup) string {

	result := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"logs:PutLogEvents",
					"logs:GetLogEvents",
					"logs:DescribeLogGroups",
					"logs:DescribeLogStreams"
				],
				"Resource": "%v:*"
			}
		]
	}`, lg.ARN)
	return result
}

// existingRole return a role with the arn value set if the role already exists
// otherwise leave the ARN value empty
func existingRole(myRole apptypes.Role) apptypes.Role {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("error loading AWS SDK configuration")
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Get IAM role input
	getRoleInput := &iam.GetRoleInput{
		RoleName: aws.String(myRole.Name),
	}

	result, err := client.GetRole(context.TODO(), getRoleInput)
	if err != nil {
		myRole.ARN = ""
		return myRole
	}

	// Role found
	myRole.ARN = *result.Role.Arn
	return myRole
}
func roleWithPolicy(
	lg *apptypes.LogGroup, cConfig appcfg.CreateConfig, log *zerolog.Logger) (r *apptypes.Role, err error) {
	result := apptypes.Role{
		Name: appcfg.RoleNameFromVPC(cConfig.VPCID),
	}
	result = existingRole(result)
	if result.ARN != "" {
		return &result, nil
	}

	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK configuration:", err)
		return
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Create IAM role input
	createRoleInput := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "vpc-flow-logs.amazonaws.com"
                    },
                    "Action": "sts:AssumeRole"
                }
            ]
        }`),
		RoleName: aws.String(result.Name),
	}

	// Create the IAM role
	createRoleOutput, err := client.CreateRole(context.TODO(), createRoleInput)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating IAM role")
		return r, err
	}
	result.ARN = *createRoleOutput.Role.Arn
	log.Debug().Msgf("IAM Role ARN: %s", result.ARN)

	// Create IAM policy input
	createPolicyInput := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(policyDocument(*lg)),
		PolicyName:     aws.String(result.Name),
	}

	// Create the IAM policy
	createPolicyOutput, err := client.CreatePolicy(context.TODO(), createPolicyInput)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating IAM policy")
		return r, err
	}

	log.Debug().Msgf("IAM Policy ARN: %s", *createPolicyOutput.Policy.Arn)

	// Attach IAM policy to IAM role
	attachPolicyInput := &iam.AttachRolePolicyInput{
		PolicyArn: createPolicyOutput.Policy.Arn,
		RoleName:  aws.String(result.Name),
	}

	_, err = client.AttachRolePolicy(context.TODO(), attachPolicyInput)
	if err != nil {
		log.Fatal().Err(err).Msg("error attaching IAM policy to IAM role")
		return r, err
	}

	log.Debug().Msgf("IAM Policy attached to IAM Role: %s", result.Name)
	r = &result
	return r, nil
}

// CreateVPCFlowLog creates a VPC Flow Log to a specified CloudWatch log group using a given IAM role
func vpcFlowLog(fl apptypes.FlowLog, log *zerolog.Logger) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("error loading AWS SDK configuration")
	}

	// Create EC2 client
	client := ec2.NewFromConfig(cfg)

	// Create VPC Flow Log input
	createFlowLogInput := &ec2.CreateFlowLogsInput{
		ResourceIds:              []string{fl.VPCID},
		ResourceType:             ec2types.FlowLogsResourceTypeVpc,
		TrafficType:              ec2types.TrafficTypeAll,
		MaxAggregationInterval:   aws.Int32(600),
		DeliverLogsPermissionArn: aws.String(fl.RoleARN),
		LogGroupName:             aws.String(fl.LogGroupName),
	}

	_, err = client.CreateFlowLogs(context.TODO(), createFlowLogInput)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating VPC Flow Log")
	}
	return nil
}

// FlowLog creates a new VPC flow log configuration
func FlowLog(cConfig appcfg.CreateConfig, log *zerolog.Logger) error {
	awsInfo, err := appcfg.GetAWSInfo()
	if err != nil {
		return fmt.Errorf("error retrieving AWS info: %v", err)
	}
	var result apptypes.FlowLog = apptypes.FlowLog{
		LogGroupName: logGroupNameFromVPC(cConfig.VPCID),
		VPCID:        cConfig.VPCID,
	}

	logGroup, err := logGroup(cConfig)
	if err != nil {
		return fmt.Errorf("error creating log group: %v", err)
	}
	result.LogGroupARN = logGroup.ARN

	lg, err := roleWithPolicy(logGroup, cConfig, log)
	if err != nil {
		return fmt.Errorf("error creating IAM role and policy: %v", err)
	}
	result.RoleName = lg.Name
	result.RoleARN = lg.ARN
	result.PolicyName = lg.Name
	result.PolicyARN = appcfg.PolicyARNFromName(lg.Name, *awsInfo)

	err = vpcFlowLog(result, log)
	if err != nil {
		return fmt.Errorf("error creating VPC Flow Log: %v", err)
	}
	fmt.Print(result.String())
	return nil
}

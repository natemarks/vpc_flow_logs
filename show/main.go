package show

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	appcfg "github.com/natemarks/vpc_flow_logs/config"
	"github.com/natemarks/vpc_flow_logs/types"
	"github.com/rs/zerolog"
)

// GetFlowLogDescriptions retrieves a slice of flow log descriptions
func GetFlowLogDescriptions(logger *zerolog.Logger) (flowLogs []types.FlowLog, err error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error loading AWS SDK configuration: %v", err)
	}
	awsInfo, err := appcfg.GetAWSInfo()
	if err != nil {
		logger.Fatal().Err(err).Msg("error retrieving AWS info")
	}
	// Create EC2 client
	client := ec2.NewFromConfig(cfg)

	// Create paginator
	paginator := ec2.NewDescribeFlowLogsPaginator(client, &ec2.DescribeFlowLogsInput{})

	// Extract flow log descriptions
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error retrieving next page of flow logs: %v", err)
		}

		for _, flowLog := range output.FlowLogs {
			logGroupName := *flowLog.LogGroupName
			logGroupArn, _ := appcfg.GetLogGroupARN(logGroupName)
			roleAndPolicyName := appcfg.RoleNameFromVPC(*flowLog.ResourceId)
			description := types.FlowLog{
				FlowLogID:      *flowLog.FlowLogId,
				VPCID:          *flowLog.ResourceId,
				RoleName:       roleAndPolicyName,
				RoleARN:        *flowLog.DeliverLogsPermissionArn,
				PolicyName:     roleAndPolicyName,
				PolicyARN:      appcfg.PolicyARNFromName(roleAndPolicyName, *awsInfo),
				LogGroupName:   logGroupName,
				LogGroupARN:    logGroupArn,
				DeliveryStatus: *flowLog.FlowLogStatus,
			}
			if err := description.Validate(); err != nil {
				logger.Error().Err(err).Msg("invalid flow log description")
				continue
			}
			flowLogs = append(flowLogs, description)
		}
	}

	return flowLogs, nil
}

// PrintFlowLogDescriptions prints the flow log descriptions to stdout
func PrintFlowLogDescriptions(config appcfg.ShowConfig, logger *zerolog.Logger) {
	fmt.Printf("Flow Log Descriptions(%v :: %v)\n", config.AccountID, config.Region)
	flowLogs, err := GetFlowLogDescriptions(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error retrieving flow log descriptions")
	}
	for _, flowLog := range flowLogs {
		fmt.Println(flowLog.String())
	}
}

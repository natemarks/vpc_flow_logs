package delete

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	appcfg "github.com/natemarks/vpc_flow_logs/config"
	apptypes "github.com/natemarks/vpc_flow_logs/types"
	"github.com/rs/zerolog"
)

// GetFlowLog retrieves a slice of flow log descriptions
func GetFlowLog(flowLogID string) (flowLog apptypes.FlowLog, err error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("error loading AWS SDK configuration: %v", err))
	}
	awsInfo, err := appcfg.GetAWSInfo()
	if err != nil {
		panic(fmt.Errorf("error retrieving AWS info: %v", err))
	}
	// Create EC2 client
	client := ec2.NewFromConfig(cfg)

	// Create paginator
	paginator := ec2.NewDescribeFlowLogsPaginator(client, &ec2.DescribeFlowLogsInput{})

	// Extract flow log descriptions
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return flowLog, fmt.Errorf("error retrieving next page of flow logs: %v", err)
		}

		for _, flowLog := range output.FlowLogs {
			if *flowLog.FlowLogId != flowLogID {
				continue
			}
			logGroupName := *flowLog.LogGroupName
			logGroupArn, _ := appcfg.GetLogGroupARN(logGroupName)
			roleAndPolicyName := appcfg.RoleNameFromVPC(*flowLog.ResourceId)
			description := apptypes.FlowLog{
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
			return description, nil
		}
	}

	return flowLog, nil
}

// FlowLog deletes a VPC Flow Log configuration
// vpc-0d354f2e35a217375_flow_log_to_cloudwatch_logs
// arn:aws:iam::151924297945:role/vpc-0d354f2e35a217375_flow_log_to_cloudwatch_logs
//
// vpc-0d354f2e35a217375_flow_log_to_cloudwatch_logs
// arn:aws:iam::151924297945:policy/vpc-0d354f2e35a217375_flow_log_to_cloudwatch_logs
func FlowLog(cfg appcfg.DeleteConfig, log *zerolog.Logger) {
	target, err := GetFlowLog(cfg.FlowLogID)
	if err != nil {
		log.Fatal().Err(err).Msg("error getting flow log")
	}
	log.Debug().Msgf("Deleting VPC Flow Log config: %v", cfg.FlowLogID)

	log.Debug().Msgf("deleting VPC flow log (%s)", target.FlowLogID)
	err = appcfg.DeleteVPCFlowLog(target.FlowLogID)
	if err != nil {
		log.Error().Err(err).Msg("error deleting VPC flow log")
	}

	log.Debug().Msgf("detaching policy (%s) from role (%s)", target.PolicyName, target.RoleName)
	err = appcfg.DetachRolePolicy(target.RoleName, target.PolicyARN)
	if err != nil {
		log.Error().Err(err).Msg("error detaching policy from role")
	}

	log.Debug().Msgf("deleting policy (%s)", target.PolicyName)
	err = appcfg.DeletePolicy(target.PolicyARN)
	if err != nil {
		log.Error().Err(err).Msg("error deleting policy")
	}

	log.Debug().Msgf("deleting role (%s)", target.RoleName)
	err = appcfg.DeleteRole(target.RoleName)
	if err != nil {
		log.Error().Err(err).Msg("error deleting role")
	}

	log.Debug().Msgf("deleting cloudwatch log group (%s)", target.LogGroupName)
	err = appcfg.DeleteCloudwatchLogGroup(target.LogGroupName)
	if err != nil {
		log.Error().Err(err).Msg("error deleting cloudwatch log group")
	}
}

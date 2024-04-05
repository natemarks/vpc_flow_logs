package types

import (
	"fmt"
	"strings"
)

// LogGroup represents a CloudWatch log group
type LogGroup struct {
	Name          string `json:"name"`
	ARN           string `json:"arn"`
	RetentionDays int    `json:"retentionDays"`
}

// Role represents an IAM role
type Role struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

// FlowLog represents a VPC flow log configuration
type FlowLog struct {
	FlowLogID      string `json:"flowLogID"`
	VPCID          string `json:"vpcID"`
	RoleName       string `json:"roleName"`
	RoleARN        string `json:"roleARN"`
	PolicyName     string `json:"policyName"`
	PolicyARN      string `json:"policyARN"`
	LogGroupName   string `json:"logGroupName"`
	LogGroupARN    string `json:"logGroupARN"`
	DeliveryStatus string `json:"deliveryStatus"`
}

// Validate checks if the flow log configuration is valid
func (f *FlowLog) Validate() (err error) {
	if f.VPCID == "" {
		return fmt.Errorf("vpc id is required")
	}
	if !strings.HasPrefix(f.VPCID, "vpc-") {
		return fmt.Errorf("vpc id must start with 'vpc-'")
	}
	if f.RoleName == "" {
		return fmt.Errorf("role name is required")
	}
	if f.LogGroupName == "" {
		return fmt.Errorf("log group name is required")
	}
	if f.LogGroupARN == "" {
		return fmt.Errorf("log group ARN is required")
	}
	return nil

}

func (f *FlowLog) String() string {
	return fmt.Sprintf("FlowLogID: %s\nVPCID: %s\nRoleName: %s\nRoleARN: %s\nPolicyName: %s\nPolicyARN: %s\nLogGroupName: %s\nLogGroupARN: %s\nDeliveryStatus: %s\n\n",
		f.FlowLogID, f.VPCID, f.RoleName, f.RoleARN, f.PolicyName, f.PolicyARN, f.LogGroupName, f.LogGroupARN, f.DeliveryStatus)
}

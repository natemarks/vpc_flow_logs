# vpc_flow_logs

It's generally not valuable to keep VPC flow logs all the time. This project
can be used to configure flow logs for cloudwatch to quickly  look at the data,
then clean up the resources.  cloud watch is a little pricey, so we set the
retention to 14 days.  

The executable has show, create and delete subcommands. 

Show will show all  the configured flow logs in the current account and region

Create will create and configure  the log groue, role, and flow log for the
specified VPC

Delete will delete the flow log, role, and log group for the specified flow log
ID


## Usage


Print help:

```bash
./build/current/linux/amd64/vpc_flow_logs
show/create/delete VPC flow logs in your account and region.

Usage:
  vpc_flow_logs [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new VPC flow log configuration
  delete      Delete a vpc flow log configuration and related resources
  help        Help about any command
  show        Print VPC flow log descriptions

Flags:
  -d, --debug   Enable debug logging
  -h, --help    help for vpc_flow_logs

Use "vpc_flow_logs [command] --help" for more information about a command.

```


Show the current flow log configurations:

```bash
./build/current/linux/amd64/vpc_flow_logs show
Flow Log Descriptions(151924297945 :: us-east-1)
FlowLogID: fl-04f78f2e41c75c80c
VPCID: vpc-0595265c52fa07048
RoleName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
RoleARN: arn:aws:iam::151924297945:role/flow_logs_to_cloudwatch_logs
PolicyName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
PolicyARN: arn:aws:iam::151924297945:policy/vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
LogGroupName: your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048:*
DeliveryStatus: ACTIVE


```

Create a new flow log configuration:

```bash
./build/current/linux/amd64/vpc_flow_logs create -v vpc-06b7636c84f111514
FlowLogID: 
VPCID: vpc-06b7636c84f111514
RoleName: vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
RoleARN: arn:aws:iam::151924297945:role/vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
PolicyName: vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
PolicyARN: arn:aws:iam::151924297945:policy/vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
LogGroupName: /aws/vpc/vpc-06b7636c84f111514/flowlogs
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:/aws/vpc/vpc-06b7636c84f111514/flowlogs
DeliveryStatus: 
```

show again to see the new configuration:

```bash
./build/current/linux/amd64/vpc_flow_logs show
Flow Log Descriptions(151924297945 :: us-east-1)
FlowLogID: fl-04f78f2e41c75c80c
VPCID: vpc-0595265c52fa07048
RoleName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
RoleARN: arn:aws:iam::151924297945:role/flow_logs_to_cloudwatch_logs
PolicyName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
PolicyARN: arn:aws:iam::151924297945:policy/vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
LogGroupName: your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048:*
DeliveryStatus: ACTIVE


FlowLogID: fl-01bbf4f7d5933983f
VPCID: vpc-06b7636c84f111514
RoleName: vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
RoleARN: arn:aws:iam::151924297945:role/vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
PolicyName: vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
PolicyARN: arn:aws:iam::151924297945:policy/vpc-06b7636c84f111514_flow_log_to_cloudwatch_logs
LogGroupName: /aws/vpc/vpc-06b7636c84f111514/flowlogs
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:/aws/vpc/vpc-06b7636c84f111514/flowlogs:*
DeliveryStatus: ACTIVE
```

Delete the flow log configuration:

```bash
./build/current/linux/amd64/vpc_flow_logs delete -f fl-01bbf4f7d5933983f
```

show again to see the new configuration:

```bash
./build/current/linux/amd64/vpc_flow_logs show
Flow Log Descriptions(151924297945 :: us-east-1)
FlowLogID: fl-04f78f2e41c75c80c
VPCID: vpc-0595265c52fa07048
RoleName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
RoleARN: arn:aws:iam::151924297945:role/flow_logs_to_cloudwatch_logs
PolicyName: vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
PolicyARN: arn:aws:iam::151924297945:policy/vpc-0595265c52fa07048_flow_log_to_cloudwatch_logs
LogGroupName: your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048:*
DeliveryStatus: ACTIVE
```

## querying the logs
go to cloudwatch logs insights and query the logs.  Here is an example query:

select your cloudwatch log and run this query for the top 10 talkers  in MB
```text
stats sum(bytes)/1048576 as megaBytesTransferred by srcAddr, dstAddr
| sort megaBytesTransferred desc
| limit 10
```

to narrow it to source and destination subnets:
```text
filter isIpv4InSubnet(dstAddr,"10.207.0.0/16") AND isIpv4InSubnet(srcAddr,"10.153.0.0/16") | stats sum(bytes)/1048576 as megaBytesTransferred by srcAddr, dstAddr
| sort megaBytesTransferred desc
| limit 10
```

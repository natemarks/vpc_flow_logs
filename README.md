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
vpc_flow_logs
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
  -h, --help     help for vpc_flow_logs
  -t, --toggle   Help message for toggle

Use "vpc_flow_logs [command] --help" for more information about a command.

```


Show the current flwo log configurations:

```bash
vpc_flow_logs show
Flow Log Descriptions(151924297945 :: us-east-1)
FlowLogID: fl-04f78f2e41c75c80c
VPCID: vpc-0595265c52fa07048
RoleName: flowlogsRole-vpc-0595265c52fa07048
RoleARN: arn:aws:iam::151924297945:role/flow_logs_to_cloudwatch_logs
LogGroupName: your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048
LogGroupARN: arn:aws:logs:us-east-1:151924297945:log-group:your-log-group-name/aws/vpc_vlow_logs/vpc-0595265c52fa07048:*
DeliveryStatus: ACTIVE

```
### Future Ideas

execute this as SSM run commands. managed the run  command documents in CDK



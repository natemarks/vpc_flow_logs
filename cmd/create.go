package cmd

// create subcommand is used to crteat a VPC flow log configuration
import (
	"os"

	"github.com/natemarks/vpc_flow_logs/config"
	"github.com/natemarks/vpc_flow_logs/create"

	"github.com/spf13/cobra"
)

var vpcID string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPC flow log configuration",
	Long: `Create a new VPC flow log configuration
	including the role, policy, log group and log stream.
	
	the log group will have a 14 day retention period`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(debug)
		cfg, err := config.GetCreateConfig(vpcID, debug)
		if err != nil {
			log.Fatal().Err(err).Msg("error getting create configuration")
			os.Exit(1)
		}
		create.FlowLog(cfg, &log)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "The VPC ID to create the flow log for")
	createCmd.MarkFlagRequired("vpc-id")
}

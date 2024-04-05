package cmd

// delete subcommand to delete a vpc flow log configuration and related resources
import (
	"os"

	"github.com/natemarks/vpc_flow_logs/config"
	"github.com/natemarks/vpc_flow_logs/delete"

	"github.com/spf13/cobra"
)

var flowLogID string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a vpc flow log configuration and related resources",
	Long: `Delete a vpc flow log configuration and related resources
	including the role, policy, log group and log stream.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(debug)
		cfg, err := config.GetDeleteConfig(flowLogID, debug)
		if err != nil {
			log.Fatal().Err(err).Msg("error getting create configuration")
			os.Exit(1)
		}
		delete.FlowLog(cfg, &log)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&flowLogID, "flow-log-id", "f", "", "the flow log ID to be deleted")
	deleteCmd.MarkFlagRequired("flow-log-id")
}

package cmd

// show subcommand
// shows the details of a vpc flow logs configs in the account and region
import (
	"os"

	"github.com/natemarks/vpc_flow_logs/config"
	"github.com/natemarks/vpc_flow_logs/show"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Print VPC flow log descriptions",
	Long:  `Print the configured VPC flow logs in the current account and region.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(false)

		cfg, err := config.GetShowConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("error getting show configuration")
			os.Exit(1)
		}
		show.PrintFlowLogDescriptions(cfg, &log)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

}

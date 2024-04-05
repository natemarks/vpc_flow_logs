package cmd

// root command for the vpc_flow_logs CLI
// creates the debug variable used by the create and delete subcommands
import (
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "vpc_flow_logs",
	Short: "View and manage VPC flow logs in your account and regi0on",
	Long:  `show/create/delete VPC flow logs in your account and region.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	// debug var shared by create and delete subcommands
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
}

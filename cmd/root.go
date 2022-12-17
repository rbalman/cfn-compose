package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configFile string
var logLevel string
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "cfnc",
	Version: "0.0.1",
	Short: "declarative way of managing CloudFormation Stacks at scale",
	Long: `Manage CloudFormation stacks at scale. Orchestrate the CloudFormation stacks just by specifying the human readable configuration. Right now yml is the only supports format`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "cfn-compose.yml", "file path to compose file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "INFO", "Specify Log Levels. Valid Levels are: DEBUG, INFO, WARN, ERROR")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "execute command in dry run mode")

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(validateCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}



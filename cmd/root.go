package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "cfnc",
	Version: "0.0.1",
	Short: "orchestrate CloudFormation stacks with ease",
	Long: `Create, Update, Delete multiple cloudformation with super ease.`,
	Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	//persistent flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "cfn-compose.yml", "config file (default is cfn-compose.yml)")
	// viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	// viper.SetDefault("config", "./cfn-compose.yml")

	//local flags
	deployCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "execute command in dry run mode")

	rootCmd.AddCommand(deployCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}



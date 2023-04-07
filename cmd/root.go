package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configFile string
var logLevel string
var dryRun bool
var flowName string
var numberOfWorkers int

var rootCmd = &cobra.Command{
	Use:     "cfnc",
	Version: "0.0.2-beta",
	Short:   "Declarative way of managing cloudformation stacks at scale",
	Long:    `Manage cloudformation stacks at scale. Design and deploy multiple cloudformation stacks either in sequence or in prallel using declarative configuration`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "cfn-compose.yml", "File path to compose file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "INFO", "Specify Log Levels. Valid Levels are: DEBUG, INFO, WARN, ERROR")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Run commands in dry run mode")
	deployCmd.PersistentFlags().StringVarP(&flowName, "flow", "f", "", "Cherry pick flow name that you want to deploy")
	deployCmd.PersistentFlags().IntVarP(&numberOfWorkers, "workers", "w", 0, "Number of concurrent workers that you want to spin")
	destroyCmd.PersistentFlags().StringVarP(&flowName, "flow", "f", "", "Cherry pick flow name that you want to destroy")
	destroyCmd.PersistentFlags().IntVarP(&numberOfWorkers, "workers", "w", 0, "Number of concurrent workers that you want to spin")

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(validateCmd)
	configCmd.AddCommand(visualizeCmd)
	configCmd.AddCommand(generateCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

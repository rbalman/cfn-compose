package cmd

import (
	"github.com/rbalman/cfn-compose/compose"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploys the stacks based on the sequence specified in the compose configuration",
	Aliases: []string{"dp"},
	Long:    `Deploys stacks based on the sequence specified in the compose configuration. Behind the scene it creates the stack if not created and updates the stack if already created. Supports dryRun mode, use --dry-run or -d flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := compose.Composer{
			LogLevel:         logLevel,
			CherryPickedFlow: flowName,
			DeployMode:       true,
			DryRun:           dryRun,
			ConfigFile:       configFile,
			NumberOfWorkers:  numberOfWorkers,
		}

		c.PrintConfig()
		c.Apply()
	},
}

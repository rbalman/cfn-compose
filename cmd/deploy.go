package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/compose"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys the stacks that are part of the compose file",
	Aliases: []string{"dp"},
	Long:  `deploy respects the order specified in the compose file and applies the changes accordingly in the individual CFN stacks. Behind the scene it creates the stack if not created and updates the stack if it already exists. Supports dryRun mode, use --dry-run or -d flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := compose.Composer{
			LogLevel: logLevel,
			CherryPickedJob: jobName,
			DeployMode: true,
			DryRun: dryRun,
			ConfigFile: configFile,
		}

		c.PrintConfig()
		c.Apply()
	},
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/compose"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroys all the stacks which are part of the compose file",
	Aliases: []string{"ds"},
	Long:  `destroy respects the order specified in the compose file and applies the changes accordingly for the individual CFN stacks. It just does the reverse thing deploy sub-command does. Supports dryRun mode, use --dry-run or -d flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := compose.Composer{
			LogLevel: logLevel,
			CherryPickedJob: jobName,
			DeployMode: false,
			DryRun: dryRun,
			ConfigFile: configFile,
		}

		c.Print()
		c.Apply()
	},
}

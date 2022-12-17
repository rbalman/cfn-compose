package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/compose"
	"path/filepath"
	"os"
)

var dryRun bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy changes in the compose file",
	Aliases: []string{"d"},
	Long:  `Creates/Updates CFN resources if required`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := filepath.Dir(configFile)
		file := filepath.Base(configFile)
		os.Chdir(dir)

		compose.Composer(file, dryRun)
	},
}

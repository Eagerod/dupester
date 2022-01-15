package dupester

import (
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dupester",
		Short: "command line tool for finding duplicates, or something",
		Args:  cobra.ArbitraryArgs,
	}
}

func Run() error {
	cmd := rootCmd()

	cmd.AddCommand(runCmd())

	return cmd.Execute()
}

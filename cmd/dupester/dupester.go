package dupester

import (
	"fmt"
)

import (
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dupester",
		Short: "command line tool for finding duplicates, or something",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Root")
			return nil
		},
	}
}

func Run() error {
	cmd := rootCmd()

	cmd.AddCommand(runCmd())

	return cmd.Execute()
}

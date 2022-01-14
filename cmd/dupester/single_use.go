package dupester

import (
	"fmt"
)

import (
	"github.com/spf13/cobra"
)

// import (
// 	"github.com/google/go-tika"
// )


func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run against a single local file",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Ran")
			return nil
		},
	}
}
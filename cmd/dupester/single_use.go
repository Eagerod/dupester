package dupester

import (
	"context"
	"fmt"
	"os"
	"strings"
)

import (
	"github.com/google/go-tika/tika"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run against a single local file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running against %s\n", args[0])

			serverUrl := "http://dev.internal.aleemhaji.com:9998"
			tikaClient := tika.NewClient(nil, serverUrl)

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			bs, err := tikaClient.ParseRecursive(context.Background(), file)
			if err != nil {
				return err
			}

			fmt.Println(strings.Join(bs, "\n"))

			return nil
		},
	}
}
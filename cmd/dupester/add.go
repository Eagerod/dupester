package dupester

import (
	"fmt"
	"path/filepath"
)

import (
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add a provided document to the pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			doc, err := dupesterClient.ParseFile(f)
			if err != nil {
				return err
			}

			existing, err := dupesterClient.FindIdentical(doc)
			if err != nil {
				return err
			}

			if existing != nil {
				return fmt.Errorf("Document identical to %s already exists", f)
			}

			return dupesterClient.Save(doc)
		},
	}
}

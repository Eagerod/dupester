package dupester

import (
	"fmt"
	"path/filepath"
)

import (
	"github.com/spf13/cobra"
)

func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "check to see if a given file has been seen before, or if there are similar files",
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
				if doc.Source == existing.Source {
					return fmt.Errorf("Document %s already seen\n", doc.Source)
				} else {
					return fmt.Errorf("Document identical to %s found at %s\n", doc.Source, existing.Source)
				}
			}

			docs, err := dupesterClient.FindLike(doc)
			if err != nil {
				return err
			}

			fmt.Printf("Found %d docs like %s\n", len(docs), f)
			if len(docs) > 0 {
				fmt.Printf("  Top match: %s\n", docs[0].Source)
			}

			return nil
		},
	}
}

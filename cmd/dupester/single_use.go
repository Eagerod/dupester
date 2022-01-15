package dupester

import (
	"fmt"
	"path/filepath"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/Eagerod/dupester/pkg/dupester"
)

type ESDoc struct {
	Source       string `json:"source"`
	OriginalBody string `json:"originalBody"`
	Body         string `json:"body"`
}

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run against a single local file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			ds, err := dupester.NewDupester("http://dev.internal.aleemhaji.com:9998", "http://dev.internal.aleemhaji.com:9200")
			if err != nil {
				return err
			}

			doc, err := ds.ParseFile(f)
			if err != nil {
				return err
			}

			existing, err := ds.FindIdentical(doc)
			if err != nil {
				return err
			}

			if existing == nil {
				docs, err := ds.FindLike(doc)
				if err != nil {
					return err
				}

				fmt.Printf("Found %d docs like %s\n", len(docs), args[0])
				if len(docs) > 0 {
					fmt.Printf("  Top match: %s\n", docs[0].Source)
				}

				err = ds.Save(doc)
				if err != nil {
					return err
				}
			} else {
				if doc.Source == existing.Source {
					fmt.Printf("Document %s already seen\n", doc.Source)
				} else {
					fmt.Printf("Document identical to %s found at %s\n", doc.Source, existing.Source)
				}
			}

			return nil
		},
	}
}

package dupester

import (
	"fmt"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/Eagerod/dupester/pkg/dupester"
)

type ESDoc struct {
	Source string `json:"source"`
	OriginalBody string `json:"originalBody"`
	Body string `json:"body"`
}

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run against a single local file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running against %s\n", args[0])

			ds, err := dupester.NewDupester("http://dev.internal.aleemhaji.com:9998", "http://dev.internal.aleemhaji.com:9200")
			if err != nil {
				return err
			}

			doc, err := ds.ParseFile(args[0])
			if err != nil {
				return err
			}

			docs, err := ds.FindLike(doc)
			if err != nil {
				return err
			}

			fmt.Printf("Found %d docs like %s\n", len(docs), args[0])

			err = ds.Save(doc)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
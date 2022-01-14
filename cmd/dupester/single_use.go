package dupester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

import (
	es "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/go-tika/tika"
	"github.com/spf13/cobra"
)

type ESDoc struct {
	Source string
	Body string
}

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run against a single local file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running against %s\n", args[0])

			rContext := context.Background()
			tikaServerUrl := "http://dev.internal.aleemhaji.com:9998"
			tikaClient := tika.NewClient(nil, tikaServerUrl)

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			bs, err := tikaClient.ParseRecursive(rContext, file)
			if err != nil {
				return err
			}

			elasticsearchServerUrl := "http://dev.internal.aleemhaji.com:9200"

			esCfg := es.Config{
				Addresses: []string{
				elasticsearchServerUrl,
				},
			  }

			elasticsearchClient, err := es.NewClient(esCfg)
			if err != nil {
				return err
			}

			o := ESDoc{args[0], strings.Join(bs, "\n")}
			b, err := json.Marshal(o)
			if err != nil {
				return err
			}

			req := esapi.IndexRequest{
				Index:      "test",
				Body:       bytes.NewReader(b),
			  }

			res, err := req.Do(rContext, elasticsearchClient)
			if err != nil {
				return err
			}

			fmt.Println(res)

			return nil
		},
	}
}
package dupester

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
    "io/ioutil"
	"os"
    "regexp"
	"strings"
)

import (
	es "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/go-tika/tika"
	"github.com/spf13/cobra"
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

			rContext := context.Background()
			tikaServerUrl := "http://dev.internal.aleemhaji.com:9998"
			tikaClient := tika.NewClient(nil, tikaServerUrl)

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			fileContents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			hash := sha256.Sum256(fileContents)

			bodyArray, err := tikaClient.ParseRecursive(rContext, file)
			if err != nil {
				return err
			}

			originalBody := strings.Join(bodyArray, "\n")

			re := regexp.MustCompile("[\\.$-/:-?{-~!\"^_`\\[\\]]")
			bodyString := re.ReplaceAllString(originalBody, " ")

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

			o := ESDoc{args[0], originalBody, bodyString}
			b, err := json.Marshal(o)
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			query := map[string]interface{}{
			  "query": map[string]interface{}{
				"more_like_this": map[string]interface{}{
				  "fields": []string{"body"},
				  "like": o.Body,
				  "min_term_freq": 1,
				  "min_doc_freq": 1,
				  "max_query_terms": 1000,
				  "analyzer": "whitespace",
				},
			  },
			}

			if err := json.NewEncoder(&buf).Encode(query); err != nil {
			  return fmt.Errorf("Error encoding query: %s", err)
			}

			// Perform the search request.
			res, err := elasticsearchClient.Search(
				elasticsearchClient.Search.WithContext(context.Background()),
				elasticsearchClient.Search.WithIndex("test"),
				elasticsearchClient.Search.WithBody(&buf),
			)

			fmt.Println(res)

			indexRequest := esapi.IndexRequest{
				DocumentID: hex.EncodeToString(hash[:]),
				Index:      "test",
				Body:       bytes.NewReader(b),
			  }

			indexResponse, err := indexRequest.Do(rContext, elasticsearchClient)
			if err != nil {
				return err
			}

			fmt.Println(indexResponse)

			return nil
		},
	}
}
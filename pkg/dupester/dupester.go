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
)

type ExtractedDocument struct {
	Source       string `json:"source"`
	Hash         string `json:"hash"`
	OriginalBody string `json:"originalBody"`
	Body         string `json:"body"`
}

type Dupester struct {
	tikaClient          *tika.Client
	elasticsearchClient *es.Client
}

func NewDupester(tikaServerUrl, elasticsearchServerUrl string) (*Dupester, error) {
	d := &Dupester{}

	d.tikaClient = tika.NewClient(nil, tikaServerUrl)

	esCfg := es.Config{
		Addresses: []string{
			elasticsearchServerUrl,
		},
	}

	elasticsearchClient, err := es.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	d.elasticsearchClient = elasticsearchClient

	return d, nil
}

func (dupester *Dupester) ParseFile(path string) (*ExtractedDocument, error) {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(fileContents)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bodyLines, err := dupester.tikaClient.ParseRecursive(context.Background(), file)
	if err != nil {
		return nil, err
	}

	originalBody := strings.Join(bodyLines, "\n")

	re := regexp.MustCompile("[\\.$-/:-?{-~!\"^_`\\[\\]]")
	bodyString := re.ReplaceAllString(originalBody, " ")

	rv := &ExtractedDocument{}
	rv.Source = path
	rv.Hash = hex.EncodeToString(hash[:])
	rv.OriginalBody = originalBody
	rv.Body = bodyString

	return rv, nil
}

func (dupester *Dupester) FindIdentical(doc *ExtractedDocument) (*ExtractedDocument, error) {
	var buf bytes.Buffer

	res, err := dupester.elasticsearchClient.Get("test", doc.Hash)
	if err != nil {
		return nil, err
	}

	buf.ReadFrom(res.Body)

	var getWrapper = struct {
		Found  bool              `json:"found"`
		Source ExtractedDocument `json:"_source"`
	}{}

	err = json.Unmarshal(buf.Bytes(), &getWrapper)
	if err != nil {
		return nil, err
	}

	if getWrapper.Found {
		return &getWrapper.Source, nil
	}
	return nil, nil
}

func (dupester *Dupester) FindLike(doc *ExtractedDocument) ([]ExtractedDocument, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"more_like_this": map[string]interface{}{
				"fields":          []string{"body"},
				"like":            doc.Body,
				"min_term_freq":   1,
				"min_doc_freq":    1,
				"max_query_terms": 1000,
				"analyzer":        "whitespace",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	res, err := dupester.elasticsearchClient.Search(
		dupester.elasticsearchClient.Search.WithContext(context.Background()),
		dupester.elasticsearchClient.Search.WithIndex("test"),
		dupester.elasticsearchClient.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, err
	}

	// Index not found, so definitely no documents found.
	rv := []ExtractedDocument{}
	if res.StatusCode == 404 {
		return rv, nil
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to query more like this, %s", res.Body)
	}

	buf.Reset()
	buf.ReadFrom(res.Body)

	var esWrapper = struct {
		Hits struct {
			Hits []struct {
				Source ExtractedDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	err = json.Unmarshal(buf.Bytes(), &esWrapper)
	if err != nil {
		return nil, nil
	}

	for _, o := range esWrapper.Hits.Hits {
		rv = append(rv, o.Source)
	}

	return rv, nil
}

func (dupester *Dupester) Save(doc *ExtractedDocument) error {
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	indexRequest := esapi.IndexRequest{
		DocumentID: doc.Hash,
		Index:      "test",
		Body:       bytes.NewReader(b),
	}

	indexResponse, err := indexRequest.Do(context.Background(), dupester.elasticsearchClient)
	if err != nil {
		return err
	}

	if indexResponse.StatusCode != 201 {
		return fmt.Errorf("Failed to create document: %s", doc.Source)
	}

	return nil
}

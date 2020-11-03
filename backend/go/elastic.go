package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticClient struct {
	Elastic *elasticsearch.Client
	Index   string
}

func NewElasticClient() *ElasticClient {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error on new elastic client: %s", err)
	}

	return &ElasticClient{Elastic: es, Index: "videogame"}
}

func (s *ElasticClient) CreateVideogame(videogame Videogame) error {
	marsh, err := json.Marshal(videogame)

	if err != nil {
		log.Printf("Error marshaling: %s", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      s.Index,
		DocumentID: videogame.Id,
		Body:       bytes.NewReader(marsh),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.Elastic)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func (s *ElasticClient) SearchVideogames(searchRequest SearchRequest) (interface{}, error) {

	return nil, errors.New("service method 'SearchVideogames' not implemented")
}

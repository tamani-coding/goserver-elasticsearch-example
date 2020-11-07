package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticClient struct {
	Elastic *elasticsearch.Client
	Index   string
}

type ElasticSearchRequest struct {
	From  int32                  `json:"from,omitempty"`
	Size  int32                  `json:"size,omitempty"`
	Sort  map[string]string      `json:"sort,omitempty"`
	Query map[string]interface{} `json:"query,omitempty"`
}

type ElasticSearchResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string    `json:"_index"`
			Type   string    `json:"_type"`
			ID     string    `json:"_id"`
			Score  float64   `json:"_score"`
			Source Videogame `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
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
	var buf bytes.Buffer

	query := ElasticSearchRequest{
		From: searchRequest.Offset,
		Size: searchRequest.Limit,
	}

	// sorting
	sort := "title.keyword"
	if searchRequest.SortBy == "publisher" {
		sort = "publisher.keyword"
	}
	if searchRequest.SortBy == "releaseDate" {
		sort = "releaseDate"
	}
	sortType := "asc"
	if searchRequest.SortType == "desc" {
		sortType = "desc"
	}
	query.Sort = map[string]string{
		sort: sortType,
	}

	// match
	match := make(map[string]interface{})
	if searchRequest.Title != "" {
		match["title"] = searchRequest.Title
	}

	// query
	query.Query = map[string]interface{}{}
	if len(match) == 0 {
		query.Query["match_all"] = make(map[string]interface{})
	} else {
		if len(match) > 0 {
			query.Query["match"] = match
		}
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s", err)
	}

	res, err := s.Elastic.Search(
		s.Elastic.Search.WithContext(context.Background()),
		s.Elastic.Search.WithIndex(s.Index),
		s.Elastic.Search.WithBody(&buf),
		s.Elastic.Search.WithTrackTotalHits(true),
		s.Elastic.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	result := ElasticSearchResult{}
	searchResult := SearchResponse{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	} else {
		log.Printf("Found %d results", result.Hits.Total.Value)
		for _, element := range result.Hits.Hits {
			searchResult.Videogames = append(searchResult.Videogames, element.Source)
		}
	}

	return searchResult, nil
}

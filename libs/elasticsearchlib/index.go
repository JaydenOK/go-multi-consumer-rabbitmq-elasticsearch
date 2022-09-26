package elasticsearchlib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"time"
)

type indexClient struct {
	es *elasticsearch.Client
}

func (i *indexClient) Close(index string) bool {

	if index == "" {
		return false
	}

	req := esapi.IndicesCloseRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return false
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
		return false
	}

	return true
}

func (i *indexClient) Exists(index string) bool {

	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return false
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
		return false
	}

	return true
}

func (i *indexClient) Create(index string, info map[string]interface{}) error {

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("index create fail, error marshaling document, index:%s, error:%s", index, err)
	}

	req := esapi.IndicesCreateRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Timeout: 30 * time.Second,
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		return fmt.Errorf("error getting response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%s] Error indexing document Index=%v", res.Status(), req.Index)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing the response body: %v", err)
		} else {
			// Print the response status and indexed document version.
			return fmt.Errorf("[%s] %s", res.Status(), r["result"])
		}
	}

	return nil
}

func (i *indexClient) Delete(indexName string) bool {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		log.Printf("Index Delete Error getting response: %s", err)
		return false
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Index Delete Error response: %s", res.String())
		return false
	}

	return true
}

func (i *indexClient) IsClose(indexName string) bool {

	req := esapi.IndicesGetSettingsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return false
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
		return false
	}

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Fatalf("Error parsing the response: %s\n", err)
	}

	if info, ok := data[indexName].(map[string]interface{}); ok {
		if settings, ok := info["settings"].(map[string]interface{}); ok {
			if index, ok := settings["index"].(map[string]interface{}); ok {
				if verifiedBeforeClose, ok := index["verified_before_close"].(string); ok {
					close, err := strconv.ParseBool(verifiedBeforeClose)
					if err != nil {
						return false
					}

					if close {
						return true
					}
				}
			}
		}
	}

	return false
}

func (i *indexClient) ForceMerge(index string) error {

	req := esapi.IndicesForcemergeRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), i.es)
	if err != nil {
		return fmt.Errorf("ForceMergeAsync error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ForceMergeAsync error response: %s", res.String())
	}

	return nil
}

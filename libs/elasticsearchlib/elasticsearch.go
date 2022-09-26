package elasticsearchlib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"strconv"
	"time"
)

var esClient *elasticsearch.Client
var bulkIndexer esutil.BulkIndexer

func InitESClient() {
	var err error
	host := viper.GetString("elasticsearch.host")
	port := viper.GetString("elasticsearch.port")
	username := viper.GetString("elasticsearch.username")
	password := viper.GetString("elasticsearch.password")
	address := "http://" + host + ":" + port
	cfg := elasticsearch.Config{
		Addresses: []string{address},
		Username:  username,
		Password:  password,
		//Transport: &http.Transport{
		//	MaxIdleConnsPerHost:   10,
		//	ResponseHeaderTimeout: time.Second,
		//	DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		//	TLSClientConfig: &tls.Config{
		//		MinVersion:         tls.VersionTLS12,
		//	},
		//},
	}
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	config := esutil.BulkIndexerConfig{
		Client:        esClient,         // The Elasticsearch client
		NumWorkers:    runtime.NumCPU(), // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	}
	bulkIndexer, err = esutil.NewBulkIndexer(config)
	if err != nil {
		log.Fatalf("Error creating BulkIndexer: %s", err)
	}
	log.Println(res)
}

func GetClient() *elasticsearch.Client {
	return esClient
}

// Index
func IndexClose(index string) bool {

	if index == "" {
		return false
	}

	req := esapi.IndicesCloseRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), esClient)
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

func IndexExists(index string) bool {

	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), esClient)
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

func IndexCreate(index string, info map[string]interface{}) error {

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("index create fail, error marshaling document, index:%s, error:%s", index, err)
	}

	req := esapi.IndicesCreateRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Timeout: 30 * time.Second,
	}

	res, err := req.Do(context.Background(), esClient)
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

func IndexDelete(indexName string) bool {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), esClient)
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

func IndexIsClose(indexName string) bool {

	req := esapi.IndicesGetSettingsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), esClient)
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

func IndexForceMerge(index string) error {

	req := esapi.IndicesForcemergeRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return fmt.Errorf("ForceMergeAsync error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ForceMergeAsync error response: %s", res.String())
	}

	return nil
}

// Document
type DocumentEntity struct {
	Id   string
	Data *map[string]interface{}
}

type Pager struct {
	pageNumber int
	pageSize   int
	totalCount int
	data       []interface{}
}

func (p *Pager) GetPageNumber() int {
	return p.pageNumber
}

func (p *Pager) GetPageSize() int {
	return p.pageSize
}

func (p *Pager) GetTotalCount() int {
	return p.totalCount
}

func (p *Pager) GetData() []interface{} {
	return p.data
}

func DocumentBatchSave(index string, docs []*DocumentEntity) error {

	for _, doc := range docs {

		data, err := json.Marshal(doc.Data)
		if err != nil {
			log.Printf("Cannot encode doc %s: %s", doc.Id, err)
		}
		err = bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Index: index,
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",
				// DocumentID is the (optional) document ID
				DocumentID: doc.Id,
				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),
				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					log.Printf("batch save success! index:%s, id:%s", res.Index, res.DocumentID)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					info, _ := json.Marshal(res)
					if err != nil {
						log.Printf("batch save has fail! info:%s ERROR: %v", info, err)
					} else {
						log.Printf("batch save has fail! info:%s ERROR: %s: %s", info, res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			return fmt.Errorf("batch save has fail! index:%s Unexpected error: %v", index, err)
		}
	}

	return nil
}

func DocumentSave(index string, doc DocumentEntity) error {

	if index == "" {
		return fmt.Errorf("document save fail, index can not be empty")
	}

	if doc.Id == "" || doc.Data == nil || len(*doc.Data) == 0 {
		return fmt.Errorf("document save fail, param doc invalid. index:%s", index)
	}

	data, err := json.Marshal(doc.Data)
	if err != nil {
		return fmt.Errorf("document save fail, error marshaling document, index:%s, error:%s", index, err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: doc.Id,
		Body:       bytes.NewReader(data),
		Timeout:    30 * time.Second,
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return fmt.Errorf("error getting response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%s] Error indexing document ID=%s, Index=%v", res.Status(), req.DocumentID, req.Index)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing the response body: %v", err)
		} else {
			// Print the response status and indexed document version.
			return fmt.Errorf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return nil
}

func DocumentFind(req esapi.SearchRequest) *Pager {

	var p = Pager{
		pageNumber: 1,
		pageSize:   20,
		totalCount: 0,
		data:       nil,
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("error getting response: %v", err)
		return &p
	}
	defer res.Body.Close()

	//判断响应码
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %v", err)
		} else {
			// Print the response status and error information.

			var errorType string
			var errorReason string
			if errorInfo, ok := e["error"].(map[string]interface{}); ok {
				if v, ok := errorInfo["type"].(string); ok {
					errorType = v
				}

				if v, ok := errorInfo["reason"].(string); ok {
					errorReason = v
				}
			}
			log.Printf("Search index fail! type:%s, reason:%s", errorType, errorReason)
		}
		return &p
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %v", err)
		return &p
	}

	if hits1, ok := r["hits"].(map[string]interface{}); ok {
		if total, ok := hits1["total"].(map[string]interface{}); ok {
			if totalCount, ok := total["value"].(float64); ok {
				p.totalCount = int(totalCount)
			}
		}

		if data, ok := hits1["hits"].([]interface{}); ok {
			p.data = data
			p.pageSize = len(data)
		}
	}

	return &p
}

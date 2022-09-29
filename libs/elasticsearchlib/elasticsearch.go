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
	"strings"
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
func IndexLists() interface{} {
	req := esapi.CatIndicesRequest{}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		fmt.Println(err.Error())
		return "查询失败:" + err.Error()
	}
	defer res.Body.Close() // res.Body  = io.ReadCloser

	buf := new(bytes.Buffer) //new(Type)作用是为T类型分配并清零一块内存，并将这块内存地址作为结果返回
	_, _ = buf.ReadFrom(res.Body)
	return strings.Split(buf.String(), "\n")
}

func IndexExist(index string) bool {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		fmt.Println("index exist err1 :", err.Error())
		return false
	}
	defer res.Body.Close()
	if res.IsError() {
		//index not exist
		return false
	}
	return true
}

// 创建Index, indexInfo 可传空map
// {"mappings":{"properties":{"name":{"type":"text"},"age":{"type":"integer"},"mobile":{"type":"text"},"company":{"type":"text"}}}}
// 文档：https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-create-index.html
func IndexCreate(index string, indexInfo map[string]interface{}) interface{} {
	data, err := json.Marshal(indexInfo)
	if err != nil {
		str := fmt.Sprintf("index create fail, error marshaling document, index:%s, error:%s", index, err)
		fmt.Println(str)
		return str
	}
	ioReader := bytes.NewReader(data)
	req := esapi.IndicesCreateRequest{
		Index:   index,
		Body:    ioReader,
		Timeout: 30 * time.Second,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("IndexCreate err : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {
		str := fmt.Sprintf("[%s] Error indexing document Index=%v, %s", res.Status(), req.Index, res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
}

// 获取mapping结构
func IndexGetMapping(index string) interface{} {
	req := esapi.IndicesGetMappingRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("req.Do : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {
		str := fmt.Sprintf("[%s] Error res， Index=%v , %s", res.Status(), req.Index, res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
}

// 增加index的mapping结构（不能修改）
// 文档：https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-put-mapping.html
func IndexPutMapping(index string, mapping map[string]interface{}) interface{} {
	byteMapping, err := json.Marshal(mapping)
	if err != nil {
		return err.Error()
	}
	reader := bytes.NewReader(byteMapping)
	req := esapi.IndicesPutMappingRequest{
		Index:             []string{index},
		Body:              reader,
		AllowNoIndices:    nil,
		ExpandWildcards:   "",
		IgnoreUnavailable: nil,
		MasterTimeout:     0,
		Timeout:           0,
		WriteIndexOnly:    nil,
		Pretty:            false,
		Human:             false,
		ErrorTrace:        false,
		FilterPath:        nil,
		Header:            nil,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("req.Do : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {
		str := fmt.Sprintf("[%s] Error res， Index=%v,%s", res.Status(), req.Index, res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
}

// 复制索引，(es不能修改直接索引名称，或修改mapping字段，需通过复制原始数据到新的索引来实现；或者添加索引别名）
func IndexReindex(sourceIndex, destIndex string) interface{} {
	data := map[string]interface{}{
		"source": map[string]string{
			"index": sourceIndex,
		},
		"dest": map[string]string{
			"index": destIndex,
		},
	}
	marshal, _ := json.Marshal(data)
	fmt.Println(string(marshal))
	body := bytes.NewReader(marshal)
	req := esapi.ReindexRequest{
		Body: body,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("req.Do : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {

		str := fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
}

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

// index别名操作（添加）
func IndexAlias(index, alias string) interface{} {
	items := make(map[string]map[string]string)
	items["add"] = map[string]string{
		"index": index,
		"alias": alias,
	}

	var actions []map[string]map[string]string
	actions = append(actions, items)
	data := map[string]interface{}{
		"actions": actions,
	}
	marshal, _ := json.Marshal(data)
	fmt.Println(string(marshal))
	body := bytes.NewReader(marshal)
	req := esapi.IndicesPutAliasRequest{
		Index: []string{index},
		Body:  body,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("req.Do : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {

		str := fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
}

func IndexAliasLists(index string) interface{} {
	req := esapi.IndicesGetAliasRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		str := fmt.Sprintf("req.Do : %v", err)
		fmt.Println(str)
		return str
	}
	defer res.Body.Close()
	if res.IsError() {

		str := fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String())
		fmt.Println(str)
		return str
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			str := fmt.Sprintf("error parsing the response body: %v", err)
			fmt.Println(str)
			return str
		}
		return r
	}
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

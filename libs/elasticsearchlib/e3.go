package elasticsearchlib

/*
import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"strings"
	"sync"
)

去弹性搜索
Elasticsearch的官方 Go 客户端。

GoDoc 去报告卡 编解码器.io 建造 单元 一体化 API

兼容性
语言客户端向前兼容；这意味着客户端支持与 Elasticsearch 更大或相同的次要版本进行通信。Elasticsearch 语言客户端仅向后兼容默认发行版，并且不做任何保证。

使用 Go 模块时，在导入路径中包含版本，并指定显式版本或分支：

require github.com/elastic/go-elasticsearch/v8 v8.0.0
require github.com/elastic/go-elasticsearch/v7 7.17
可以在单个项目中使用多个版本的客户端：

// go.mod
github.com/elastic/go-elasticsearch/v7 v7.17.0
github.com/elastic/go-elasticsearch/v8 v8.0.0

// main.go
import (
elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)
// ...
es7, _ := elasticsearch7.NewDefaultClient()
es8, _ := elasticsearch8.NewDefaultClient()
客户端的分支与 Elasticsearchmain的当前分支兼容。master

安装
将包添加到您的go.mod文件中：

require github.com/elastic/go-elasticsearch/v8 main
或者，克隆存储库：

git clone --branch main https://github.com/elastic/go-elasticsearch.git $GOPATH/src/github.com/elastic/go-elasticsearch
一个完整的例子：

mkdir my-elasticsearch-app && cd my-elasticsearch-app

cat > go.mod <<-END
module my-elasticsearch-app

require github.com/elastic/go-elasticsearch/v8 main
END

cat > main.go <<-END
package main

import (
"log"

"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
END

go run main.go
用法
该elasticsearch包将两个单独的包联系在一起，分别用于调用 Elasticsearch API 和通过 HTTP:esapi和传输数据elastictransport。

使用该elasticsearch.NewDefaultClient()功能以默认设置创建客户端。

es, err := elasticsearch.NewDefaultClient()
if err != nil {
log.Fatalf("Error creating the client: %s", err)
}

res, err := es.Info()
if err != nil {
log.Fatalf("Error getting response: %s", err)
}

defer res.Body.Close()
log.Println(res)

// [200 OK] {
//   "name" : "node-1",
//   "cluster_name" : "go-elasticsearch"
// ...
注意：关闭响应正文并使用它至关重要，以便在默认 HTTP 传输中重用持久 TCP 连接。如果您对响应正文不感兴趣，请致电.io.Copy(ioutil.Discard, res.Body)

导出ELASTICSEARCH_URL环境变量时，它将用于设置集群端点。用逗号分隔多个地址。

要以编程方式设置集群端点，请将配置对象传递给elasticsearch.NewClient()函数。

cfg := elasticsearch.Config{
Addresses: []string{
"https://localhost:9200",
"https://localhost:9201",
},
// ...
}
es, err := elasticsearch.NewClient(cfg)
要设置用户名和密码，请将它们包含在端点 URL 中，或使用相应的配置选项。

cfg := elasticsearch.Config{
// ...
Username: "foo",
Password: "bar",
}
要设置用于签署集群节点证书的自定义证书颁发机构，请使用CACert配置选项。

cert, _ := ioutil.ReadFile(*cacert)

cfg := elasticsearch.Config{
// ...
CACert: cert,
}
要设置指纹以验证 HTTPS 连接，请使用CertificateFingerprint配置选项。

cfg := elasticsearch.Config{
// ...
CertificateFingerprint: fingerPrint,
}
要配置其他 HTTP 设置，http.Transport 请在配置对象中传递一个对象。

cfg := elasticsearch.Config{
Transport: &http.Transport{
MaxIdleConnsPerHost:   10,
ResponseHeaderTimeout: time.Second,
TLSClientConfig: &tls.Config{
MinVersion: tls.VersionTLS12,
// ...
},
// ...
},
}
有关客户端配置和自定义的更多示例，请参阅_examples/configuration.go和 _examples/customization.go文件。_examples/security有关安全配置的示例，请参见。

以下示例演示了更复杂的用法。它从集群中获取 Elasticsearch 版本，同时索引几个文档，并使用响应正文周围的轻量级包装器打印搜索结果。

// $ go run _examples/main.go

package main

import (
"bytes"
"context"
"encoding/json"
"log"
"strconv"
"strings"
"sync"

"github.com/elastic/go-elasticsearch/v8"
"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	log.SetFlags(0)

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// 2. Index documents concurrently
	//
	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			// Build the request body.
			data, err := json.Marshal(struct{ Title string }{Title: title})
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))

	// 3. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}

// Client: 8.0.0-SNAPSHOT
// Server: 8.0.0-SNAPSHOT
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// [201 Created] updated; version=1
// [201 Created] updated; version=1
// -------------------------------------
// [200 OK] 2 hits; took: 5ms
//  * ID=1, map[title:Test One]
//  * ID=2, map[title:Test Two]
// =====================================
正如您在上面的示例中看到的，该esapi包允许以两种不同的方式调用 Elasticsearch API：或者通过创建一个结构，例如IndexRequest，并通过向其传递上下文和客户端来调用其方法Do()，或者通过调用Search()客户端直接，使用选项功能，如WithIndex(). 在包文档中查看更多信息和示例 。

该elastictransport包处理与 Elasticsearch 之间的数据传输，包括重试失败的请求、保持连接池、发现集群节点和日志记录。

在以下博客文章中阅读有关客户端内部结构和使用的更多信息：

https://www.elastic.co/blog/the-go-client-for-elasticsearch-introduction
https://www.elastic.co/blog/the-go-client-for-elasticsearch-configuration-and-customization
https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data
帮手
该esutil软件包提供了与客户合作的便利助手。目前，它提供了 theesutil.JSONReader()和esutil.BulkIndexerhelpers。

例子
该_examples文件夹包含许多配方和综合示例，可帮助您开始使用客户端，包括客户端的配置和自定义、使用自定义证书颁发机构 (CA) 以确保安全 (TLS)、模拟传输以进行单元测试、嵌入客户端在自定义类型中，构建查询，单独和批量执行请求，并解析响应。

执照
该软件在Apache 2 许可下获得许可。见通知。
*/

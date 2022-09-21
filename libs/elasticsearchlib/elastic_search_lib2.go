package elasticsearchlib

//
//import (
//	"fmt"
//	"log"
//	"os"
//	"reflect"
//)
//
//ElasticSearch介绍
//1.1.1. 介绍
//Elasticsearch（ES）是一个基于Lucene构建的开源、分布式、RESTful接口的全文搜索引擎。Elasticsearch还是一个分布式文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索，ES能够横向扩展至数以百计的服务器存储以及处理PB级的数据。可以在极短的时间内存储、搜索和分析大量的数据。通常作为具有复杂搜索场景情况下的核心发动机。
//
//1.1.2. Elasticsearch能做什么
//当你经营一家网上商店，你可以让你的客户搜索你卖的商品。在这种情况下，你可以使用ElasticSearch来存储你的整个产品目录和库存信息，为客户提供精准搜索，可以为客户推荐相关商品。
//当你想收集日志或者交易数据的时候，需要分析和挖掘这些数据，寻找趋势，进行统计，总结，或发现异常。在这种情况下，你可以使用Logstash或者其他工具来进行收集数据，当这引起数据存储到ElasticsSearch中。你可以搜索和汇总这些数据，找到任何你感兴趣的信息。
//对于程序员来说，比较有名的案例是GitHub，GitHub的搜索是基于ElasticSearch构建的，在http://github.com/search页面，你可以搜索项目、用户、issue、pull request，还有代码。共有40~50个索引库，分别用于索引网站需要跟踪的各种数据。虽然只索引项目的主分支（master），但这个数据量依然巨大，包括20亿个索引文档，30TB的索引文件。
//1.1.3. Elasticsearch基本概念
//Near Realtime(NRT) 几乎实时
//Elasticsearch是一个几乎实时的搜索平台。意思是，从索引一个文档到这个文档可被搜索只需要一点点的延迟，这个时间一般为毫秒级。
//
//Cluster 集群
//群集是一个或多个节点（服务器）的集合， 这些节点共同保存整个数据，并在所有节点上提供联合索引和搜索功能。一个集群由一个唯一集群ID确定，并指定一个集群名（默认为“elasticsearch”）。该集群名非常重要，因为节点可以通过这个集群名加入群集，一个节点只能是群集的一部分。
//
//确保在不同的环境中不要使用相同的群集名称，否则可能会导致连接错误的群集节点。例如，你可以使用logging-dev、logging-stage、logging-prod分别为开发、阶段产品、生产集群做记录。
//
//Node节点
//节点是单个服务器实例，它是群集的一部分，可以存储数据，并参与群集的索引和搜索功能。就像一个集群，节点的名称默认为一个随机的通用唯一标识符（UUID），确定在启动时分配给该节点。如果不希望默认，可以定义任何节点名。这个名字对管理很重要，目的是要确定你的网络服务器对应于你的ElasticSearch群集节点。
//
//我们可以通过群集名配置节点以连接特定的群集。默认情况下，每个节点设置加入名为“elasticSearch”的集群。这意味着如果你启动多个节点在网络上，假设他们能发现彼此都会自动形成和加入一个名为“elasticsearch”的集群。
//
//在单个群集中，你可以拥有尽可能多的节点。此外，如果“elasticsearch”在同一个网络中，没有其他节点正在运行，从单个节点的默认情况下会形成一个新的单节点名为”elasticsearch”的集群。
//
//Index索引
//索引是具有相似特性的文档集合。例如，可以为客户数据提供索引，为产品目录建立另一个索引，以及为订单数据建立另一个索引。索引由名称（必须全部为小写）标识，该名称用于在对其中的文档执行索引、搜索、更新和删除操作时引用索引。在单个群集中，你可以定义尽可能多的索引。
//
//Type类型
//在索引中，可以定义一个或多个类型。类型是索引的逻辑类别/分区，其语义完全取决于你。一般来说，类型定义为具有公共字段集的文档。例如，假设你运行一个博客平台，并将所有数据存储在一个索引中。在这个索引中，你可以为用户数据定义一种类型，为博客数据定义另一种类型，以及为注释数据定义另一类型。
//
//Document文档
//文档是可以被索引的信息的基本单位。例如，你可以为单个客户提供一个文档，单个产品提供另一个文档，以及单个订单提供另一个文档。本文件的表示形式为JSON（JavaScript Object Notation）格式，这是一种非常普遍的互联网数据交换格式。
//
//在索引/类型中，你可以存储尽可能多的文档。请注意，尽管文档物理驻留在索引中，文档实际上必须索引或分配到索引中的类型。
//
//Shards & Replicas分片与副本
//索引可以存储大量的数据，这些数据可能超过单个节点的硬件限制。例如，十亿个文件占用磁盘空间1TB的单指标可能不适合对单个节点的磁盘或可能太慢服务仅从单个节点的搜索请求。
//
//为了解决这一问题，Elasticsearch提供细分你的指标分成多个块称为分片的能力。当你创建一个索引，你可以简单地定义你想要的分片数量。每个分片本身是一个全功能的、独立的“指数”，可以托管在集群中的任何节点。
//
//Shards分片的重要性主要体现在以下两个特征：
//
//1.副本为分片或节点失败提供了高可用性。为此，需要注意的是，一个副本的分片不会分配在同一个节点作为原始的或主分片，副本是从主分片那里复制过来的。
//
//2.副本允许用户扩展你的搜索量或吞吐量，因为搜索可以在所有副本上并行执行。
//
//ES基本概念与关系型数据库的比较
//ES概念	关系型数据库
//Index（索引）支持全文检索	Database（数据库）
//Type（类型）	Table（表）
//Document（文档），不同文档可以有不同的字段集合	Row（数据行）
//Field（字段）	Column（数据列）
//Mapping（映射）	Schema（模式）
//1.1.4. ES API
//以下示例使用curl演示。
//
//查看健康状态
//curl -X GET 127.0.0.1:9200/_cat/health?v
//输出：
//
//epoch      timestamp cluster       status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
//1564726309 06:11:49  elasticsearch yellow          1         1      3   3    0    0        1             0                  -                 75.0%
//查询当前es集群中所有的indices
//curl -X GET 127.0.0.1:9200/_cat/indices?v
//输出：
//
//health status index                uuid                   pri rep docs.count docs.deleted store.size pri.store.size
//green  open   .kibana_task_manager LUo-IxjDQdWeAbR-SYuYvQ   1   0          2            0     45.5kb         45.5kb
//green  open   .kibana_1            PLvyZV1bRDWex05xkOrNNg   1   0          4            1     23.9kb         23.9kb
//yellow open   user                 o42mIpDeSgSWZ6eARWUfKw   1   1          0            0       283b           283b
//创建索引
//curl -X PUT 127.0.0.1:9200/www
//输出：
//
//{"acknowledged":true,"shards_acknowledged":true,"index":"www"}
//删除索引
//curl -X DELETE 127.0.0.1:9200/www
//输出：
//
//{"acknowledged":true}
//插入记录
//curl -H "ContentType:application/json" -X POST 127.0.0.1:9200/user/person -d '
//    {
//"name": "LMH",
//"age": 18,
//"married": true
//}'
//输出：
//
//{
//"_index": "user",
//"_type": "person",
//"_id": "MLcwUWwBvEa8j5UrLZj4",
//"_version": 1,
//"result": "created",
//"_shards": {
//"total": 2,
//"successful": 1,
//"failed": 0
//},
//"_seq_no": 3,
//"_primary_term": 1
//}
//也可以使用PUT方法，但是需要传入id
//
//curl -H "ContentType:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
//    {
//"name": "LMH",
//"age": 18,
//"married": false
//}'
//检索
//Elasticsearch的检索语法比较特别，使用GET方法携带JSON格式的查询条件。
//
//全检索：
//
//curl -X GET 127.0.0.1:9200/user/person/_search
//按条件检索：
//
//curl -H "ContentType:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
//    {
//"query":{
//"match": {"name": "LMH"}
//}
//}'
//ElasticSearch默认一次最多返回10条结果，可以像下面的示例通过size字段来设置返回结果的数目。
//
//curl -H "ContentType:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
//    {
//"query":{
//"match": {"name": "LMH"},
//"size": 2
//}
//}'
//Elasticsearch安装
//Elasticsearch官网：https://www.elastic.co/cn/products/elasticsearch
//
//1.1.1. Elasticsearch介绍
//Elasticsearch（ES）是一个基于Lucene构建的开源、分布式、RESTful接口的全文搜索引擎。Elasticsearch还是一个分布式文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索，ES能够横向扩展至数以百计的服务器存储以及处理PB级的数据。可以在极短的时间内存储、搜索和分析大量的数据。通常作为具有复杂搜索场景情况下的核心发动机。
//
//
//
//
//
//
//1.1.2. 下载
//官方网站下载链接：https://www.elastic.co/cn/downloads/elasticsearch
//
//
//
//
//
//
//请根据自己的需求下载对应的版本。
//
//1.1.3. 安装
//将上一步下载的压缩包解压，下图以Windows为例。
//
//
//
//
//
//
//1.1.4. 启动
//执行bin\elasticsearch.bat启动，默认在本机的9200端口启动服务。
//
//使用浏览器访问elasticsearch服务，可以看到类似下面的信息。
//
//
//
//
//Kibana安装
//1.1.1. Kibana介绍
//官网链接：https://www.elastic.co/cn/products/kibana
//
//Kibana是一个开源的分析和可视化平台，设计用于和Elasticsearch一起工作。
//
//你可以使用Kibana来搜索、查看、并和存储在Elasticsearch索引中的数据进行交互。
//
//你可以轻松地执行高级数据分析，并且以各种图标、表格和地图的形式可视化数据。
//
//Kibana使得理解大量数据变得很容易。它简单的、基于浏览器的界面使你能够快速创建和共享动态仪表板，实时显示Elasticsearch查询的变化。
//
//1.1.2. 下载
//官方下载链接：https://www.elastic.co/cn/downloads/kibana
//
//请根据需求下载对应的版本。
//
//注意 :
//
//Kibana与Elasticsearch的版本要相互对应，否则可能不兼容！！！
//
//例如：Elasticsearch是7.2.1的版本，那么你的Kibana也要下载7.2.1的版本。
//
//
//
//
//
//
//1.1.3. 安装
//将上一步下载得到的文件解压。
//
//
//
//
//
//
//修改config目录下的配置文件kibana.yml（如你是本机没有发生改变可以省略这一步）
//
//将配置文件中 elasticsearch.hosts设置为你elasticseatch的地址，例如：
//
//
//
//
//
//
//​（找不到直接Ctrl+F搜索‘url’）
//
//然后翻到最后修改一下语言，配置成简体中文。
//
//
//
//
//
//
//1.1.4. 启动
//执行bin\kibana.bat启动
//
//启动过程比较慢，请耐心等待出现类似下图界面，就表示启动成功。
//
//
//
//
//
//
//使用浏览器访问本机的5601端口即可看到类似下面的界面：
//
//
//
//
//操作ElasticSearch
//1.1.1. elastic client
//我们使用第三方库https://github.com/olivere/elastic 来连接ES并进行操作。
//
//注意下载与你的ES相同版本的client，例如我们这里使用的ES是7.2.1的版本，那么我们下载的client也要与之对应为http://github.com/olivere/elastic/v7。
//
//使用go.mod来管理依赖：
//
//require (
//github.com/olivere/elastic/v7 v7.0.4
//)
//简单示例：
//
//package main
//
//import (
//"context"
//"fmt"
//
//"github.com/olivere/elastic/v7"
//)
//
//// Elasticsearch demo
//
//type Person struct {
//	Name    string `json:"name"`
//	Age     int    `json:"age"`
//	Married bool   `json:"married"`
//}
//
//func main() {
//	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
//	if err != nil {
//		// Handle error
//		panic(err)
//	}
//
//	fmt.Println("connect to es success")
//	p1 := Person{Name: "lmh", Age: 18, Married: false}
//	put1, err := client.Index().
//		Index("user").
//		BodyJson(p1).
//		Do(context.Background())
//	if err != nil {
//		// Handle error
//		panic(err)
//	}
//	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
//}
//示例2：
//
//package main
//
//import (
//"context"
//"fmt"
//"log"
//"os"
//"reflect"
//
//"gopkg.in/olivere/elastic.v7" //这里使用的是版本5，最新的是6，有改动
//)
//
//var client *elastic.Client
//var host = "http://127.0.0.1:9200/"
//
//type Employee struct {
//	FirstName string   `json:"first_name"`
//	LastName  string   `json:"last_name"`
//	Age       int      `json:"age"`
//	About     string   `json:"about"`
//	Interests []string `json:"interests"`
//}
//
////初始化
//func init() {
//	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
//	var err error
//	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
//	if err != nil {
//		panic(err)
//	}
//	info, code, err := client.Ping(host).Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
//
//	esversion, err := client.ElasticsearchVersion(host)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Elasticsearch version %s\n", esversion)
//
//}
//
///*下面是简单的CURD*/
//
////创建
//func create() {
//
//	//使用结构体
//	e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
//	put1, err := client.Index().
//		Index("megacorp").
//		Type("employee").
//		Id("1").
//		BodyJson(e1).
//		Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)
//
//	//使用字符串
//	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
//	put2, err := client.Index().
//		Index("megacorp").
//		Type("employee").
//		Id("2").
//		BodyJson(e2).
//		Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)
//
//	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
//	put3, err := client.Index().
//		Index("megacorp").
//		Type("employee").
//		Id("3").
//		BodyJson(e3).
//		Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)
//
//}
//
////删除
//func delete() {
//
//	res, err := client.Delete().Index("megacorp").
//		Type("employee").
//		Id("1").
//		Do(context.Background())
//	if err != nil {
//		println(err.Error())
//		return
//	}
//	fmt.Printf("delete result %s\n", res.Result)
//}
//
////修改
//func update() {
//	res, err := client.Update().
//		Index("megacorp").
//		Type("employee").
//		Id("2").
//		Doc(map[string]interface{}{"age": 88}).
//		Do(context.Background())
//	if err != nil {
//		println(err.Error())
//	}
//	fmt.Printf("update age %s\n", res.Result)
//
//}
//
////查找
//func gets() {
//	//通过id查找
//	get1, err := client.Get().Index("megacorp").Type("employee").Id("2").Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	if get1.Found {
//		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
//	}
//}
//
////搜索
//func query() {
//	var res *elastic.SearchResult
//	var err error
//	//取所有
//	res, err = client.Search("megacorp").Type("employee").Do(context.Background())
//	printEmployee(res, err)
//
//	//字段相等
//	q := elastic.NewQueryStringQuery("last_name:Smith")
//	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
//	if err != nil {
//		println(err.Error())
//	}
//	printEmployee(res, err)
//
//	//条件查询
//	//年龄大于30岁的
//	boolQ := elastic.NewBoolQuery()
//	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
//	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
//	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
//	printEmployee(res, err)
//
//	//短语搜索 搜索about字段中有 rock climbing
//	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
//	res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
//	printEmployee(res, err)
//
//	//分析 interests
//	aggs := elastic.NewTermsAggregation().Field("interests")
//	res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
//	printEmployee(res, err)
//
//}
//
////简单分页
//func list(size, page int) {
//	if size < 0 || page < 1 {
//		fmt.Printf("param error")
//		return
//	}
//	res, err := client.Search("megacorp").
//		Type("employee").
//		Size(size).
//		From((page - 1) * size).
//		Do(context.Background())
//	printEmployee(res, err)
//
//}
//
////打印查询到的Employee
//func printEmployee(res *elastic.SearchResult, err error) {
//	if err != nil {
//		print(err.Error())
//		return
//	}
//	var typ Employee
//	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
//		t := item.(Employee)
//		fmt.Printf("%#v\n", t)
//	}
//}
//
//func main() {
//	create()
//	delete()
//	update()
//	gets()
//	query()
//	list(1, 3)
//}
//更多使用详见文档：https://godoc.org/github.com/olivere/elastic
//
//发布于 2021-01-28 11:17
//elastic search
//Go 语言
//Go 编程
//​赞同 3​
//​添加评论
//​分享
//​喜欢
//​收藏
//​申请转载
//​
//写下你的评论...
//
//还没有评论，发表第一个评论吧
//文章被以下专栏收录
//Go 并发编程
//Go 并发编程
//Go 并发编程 学习与交流心得
//推荐阅读
//python 操作 ElasticSearch 入门
//python 操作 ElasticSearch 入门
//王书成
//降维打击！使用ElasticSearch作为时序数据库
//降维打击！使用ElasticSearch作为时序数据库
//Golio...
//发表于玩转Ela...
//ElasticSearch 索引 VS MySQL 索引
//前言这段时间在维护产品的搜索功能，每次在管理台看到 elasticsearch 这么高效的查询效率我都很好奇他是如何做到的。这甚至比在我本地使用 MySQL 通过主键的查询速度还快。为此我搜索了相关…
//
//cross...
//发表于cross...
//DSL不好用？一文教你如何使用ElasticSearch的内置原生SQL查询
//DSL不好用？一文教你如何使用ElasticSearch的内置原生SQL查询
//python大大
//
//
//选择语言
//登录即可查看 超5亿 专业优质内容
//超 5 千万创作者的优质提问、专业回答、深度文章和精彩视频尽在知乎。
//立即登录/注册

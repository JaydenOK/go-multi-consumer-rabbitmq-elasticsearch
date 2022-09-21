package elasticsearchlib

//
//import (
//	"fmt"
//	"reflect"
//	"strings"
//)
//
//
//package main
//
//import (
//"context"
//"encoding/json"
//"fmt"
//"github.com/olivere/elastic"
//log "github.com/sirupsen/logrus"
//"reflect"
//"strings"
//)
//
//var ESClient *elastic.Client
//var ESServerURL = []string{"http://127.0.0.1:9200"}
//
//type Language struct {
//	Name      string `json:"name"`
//	BuildTime int    `json:"build_time"`
//}
//
//func main() {
//	var esIndex = "programming"
//	var esType = "language"
//	//初始化es连接
//	ESClient, err := elastic.NewClient(
//		elastic.SetSniff(false),
//		elastic.SetURL(ESServerURL...))
//	if err != nil {
//		log.Errorf("Failed to build elasticsearch connection: %s %s", strings.Join(ESServerURL, ","), err.Error())
//	}
//	info, code, err := ESClient.Ping(strings.Join(ESServerURL, ",")).Do(context.Background())
//	if err != nil {
//		log.Error("ping es failed", err.Error())
//	}
//	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
//
//	//构造数据
//	c := &Language{
//		Name:      "c",
//		BuildTime: 1972,
//	}
//	php := &Language{
//		Name:      "php",
//		BuildTime: 1995,
//	}
//	java := Language{
//		Name:      "java",
//		BuildTime: 1995,
//	}
//	python := Language{
//		Name:      "python",
//		BuildTime: 1991,
//	}
//
//	//检查索引是否存在
//	exists, err := ESClient.IndexExists(esIndex).Do(context.Background())
//	if err != nil {
//		log.Error("check index exist failed", err.Error())
//	}
//	//索引不存在则创建索引
//	//索引不存在时查询会报错，但索引不存在的时候可以直接插入
//	if !exists {
//		log.Infof("index %s is not exist", esIndex)
//		_, err := ESClient.CreateIndex(esIndex).Do(context.Background())
//		if err != nil {
//			log.Error("create index failed", err.Error())
//		}
//	}
//
//	//插入一条数据，id指定为对应的name，若不指定则随机生成
//	_, err = ESClient.Index().Index(esIndex).Type(esType).Id(c.Name).BodyJson(c).Do(context.Background())
//	if err != nil {
//		log.Error("insert es failed", err.Error())
//	}
//
//	//借助bulk批量插入数据
//	bulkRequest := ESClient.Bulk().Index(esIndex).Type(esType)
//	req := elastic.NewBulkIndexRequest().Doc(java)
//	req.Id(java.Name)       //指定id
//	bulkRequest.Add(req)
//	req = elastic.NewBulkIndexRequest().Doc(php)
//	req.Id(php.Name)
//	bulkRequest.Add(req)
//	req = elastic.NewBulkIndexRequest().Doc(python)
//	req.Id(python.Name)
//	bulkRequest.Add(req)
//	_, err = bulkRequest.Do(context.TODO())
//	if err != nil {
//		log.Error(err.Error())
//	}
//
//	//更新数据,DetectNoop(false)时无论更新内容与原本内容是否一致，都进行更新
//	c.BuildTime = 2020
//	_,err = ESClient.Update().Index(esIndex).Type(esType).Id("c").Doc(c).DetectNoop(false).Do(context.Background())
//	if err != nil {
//		log.Error(err.Error())
//	}
//
//	//根据id进行查询
//	var resultType Language
//	searchById,err := ESClient.Get().Index(esIndex).Type(esType).Id("java").Do(context.Background())
//	if searchById.Found{
//
//		if err := json.Unmarshal(searchById.Source,&resultType); err != nil{
//			log.Error(err.Error())
//		}
//		fmt.Printf("search by id: %#v \n",resultType)
//	}
//
//	//查询index中所有的数据
//	searchAll,err := ESClient.Search(esIndex).Type(esType).Do(context.Background())
//	for _,item := range searchAll.Each(reflect.TypeOf(resultType)) {
//		language := item.(Language)
//		fmt.Printf("search by index all: %#v \n",language)
//	}
//
//	//查询前size条数据
//	//利用size和from可以实现查询结果的分页
//	searchPart,err := ESClient.Search(esIndex).Type(esType).Size(2).Do(context.Background())
//	for _,item := range searchPart.Each(reflect.TypeOf(resultType)) {
//		language := item.(Language)
//		fmt.Printf("search by index part: %#v \n",language)
//	}
//
//	//boolquery 可用于组合查询
//	//Must想当于且，Should相当于或，MustNot相当于非......
//	boolquery := elastic.NewBoolQuery()
//	boolquery.Must(elastic.NewMatchQuery("name","java"))    //查询name为java的
//	searchByMatch,err := ESClient.Search(esIndex).Type(esType).Query(boolquery).Do(context.Background())
//	for _,item := range searchByMatch.Each(reflect.TypeOf(resultType)) {
//		language := item.(Language)
//		fmt.Printf("search by match: %#v \n",language)
//	}
//
//	//匹配查询
//	matchPhraseQuery := elastic.NewMatchPhraseQuery("name","py")     //查询name包含py的
//	searchByPhrase,err := ESClient.Search(esIndex).Type(esType).Query(matchPhraseQuery).Do(context.Background())
//	for _,item := range searchByPhrase.Each(reflect.TypeOf(resultType)) {
//		language := item.(Language)
//		fmt.Printf("search by phrase: %#v \n",language)
//	}
//
//	//条件查询
//	boolquery2 := elastic.NewBoolQuery()
//	boolquery2.Filter(elastic.NewRangeQuery("build_time").Gt(2000))     //查询build_time大于2000
//	searchByfilter,err := ESClient.Search(esIndex).Type(esType).Query(boolquery2).Do(context.Background())
//	for _,item := range searchByfilter.Each(reflect.TypeOf(resultType)) {
//		language := item.(Language)
//		fmt.Printf("search by filter: %#v \n",language)
//	}
//
//	//删除指定id对应的数据
//	_,err = ESClient.Delete().Index(esIndex).Type(esType).Id("c").Do(context.Background())
//	if err != nil {
//		log.Error(err.Error())
//	}
//
//	//删除指定index中的所有数据
//	_, err = ESClient.DeleteIndex(esIndex).Do(context.Background())
//	if err != nil {
//		log.Error(err.Error())
//	}
//}

package elasticsearchlib

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
)

var esClient *elasticsearch.Client

func InitESClient() {
	var err error
	host := viper.GetString("elasticsearch.host")
	port := viper.GetString("elasticsearch.port")
	username := viper.GetString("elasticsearch.username")
	password := viper.GetString("elasticsearch.password")
	address := "http://" + host + ":" + port
	cfg := elasticsearch.Config{
		Addresses:               []string{address},
		Username:                username,
		Password:                password,
		CloudID:                 "",
		APIKey:                  "",
		ServiceToken:            "",
		CertificateFingerprint:  "",
		Header:                  nil,
		CACert:                  nil,
		RetryOnStatus:           nil,
		DisableRetry:            false,
		MaxRetries:              0,
		RetryOnError:            nil,
		CompressRequestBody:     false,
		DiscoverNodesOnStart:    false,
		DiscoverNodesInterval:   0,
		EnableMetrics:           false,
		EnableDebugLogger:       false,
		EnableCompatibilityMode: false,
		DisableMetaHeader:       false,
		RetryBackoff:            nil,
		Transport:               nil,
		Logger:                  nil,
		Selector:                nil,
		ConnectionPoolFunc:      nil,
	}
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
}

func GetClient() *elasticsearch.Client {
	return esClient
}

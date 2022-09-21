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
	log.Println(res)
}

func GetClient() *elasticsearch.Client {
	return esClient
}

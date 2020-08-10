package es

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

var client *elastic.Client
var host = "http://127.0.0.1:9200/"

func init() {
	errorLog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	client, err = elastic.NewClient(elastic.SetErrorLog(errorLog), elastic.SetURL(host))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esVersion, err := client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esVersion)
}

// Es查询
func Search(indexName string, source map[string]interface{}, from int, size int, sortByFields []elastic.Sorter) *elastic.SearchResult {
	searchResult, err := client.Search().Index(indexName).Source(source).From(from).Size(size).SortBy(sortByFields...).Pretty(true).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	return searchResult
}

// 添加ES
func Add(indexName string, id string, body map[string]interface{}) elastic.IndexResponse {
	put, err := client.Index().Index(indexName).Id(id).BodyJson(body).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return *put
}

// 修改es
func Update(indexName string, id string, body map[string]interface{}) elastic.UpdateResponse {
	update, err := client.Update().Index(indexName).Id(id).Doc(body).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return *update
}

// 删除es
func Delete(indexName string, id string) elastic.DeleteResponse {
	deleteResponse, err := client.Delete().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return *deleteResponse
}

/// 查询es

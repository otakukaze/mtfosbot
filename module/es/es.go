package es

import (
	"context"
	"fmt"

	"git.trj.tw/golang/mtfosbot/module/config"
	"github.com/olivere/elastic"
)

var client *elastic.Client

// NewClient -
func NewClient() (err error) {
	conf := config.GetConf()
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(conf.Elasticsearch.Host))
	fmt.Println("host ", conf.Elasticsearch.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

// PutLog -
func PutLog(t string, body map[string]interface{}) (err error) {
	if len(t) == 0 || body == nil {
		return
	}
	conf := config.GetConf()
	ctx := context.Background()
	_, err = client.Index().Index(conf.Elasticsearch.Index).Type(t).BodyJson(body).Do(ctx)
	return
}

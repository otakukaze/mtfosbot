package es

import (
	"context"
	"fmt"
	"time"

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
	if client == nil {
		return
	}
	if len(t) == 0 || body == nil {
		return
	}
	conf := config.GetConf()
	ctx := context.Background()
	body["timestamp"] = time.Now().UnixNano() / 1000000
	body["date"] = time.Now().UTC()
	nt := time.Now()
	idx := fmt.Sprintf("%v-%v-%v-%v", conf.Elasticsearch.Index, nt.Year(), nt.Month(), nt.Day())
	_, err = client.Index().Index(idx).Type(t).BodyJson(body).Do(ctx)
	return
}

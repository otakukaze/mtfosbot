package es

import (
  "context"
  "fmt"
  "git.trj.tw/golang/mtfosbot/module/config"
  "github.com/olivere/elastic"
)

var client *elastic.Client

func NewClient () (err error) {
  conf := config.GetConf()
  client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(conf.Elasticsearch.Host))
  fmt.Println("host ", conf.Elasticsearch.Host)
  if err != nil {
    fmt.Println(err)
    return
  }
  return
}

func PutLog () (err error) {
  conf := config.GetConf()
  ctx := context.Background()
  _, err = client.Index().Index(conf.Elasticsearch.Index).Type("type").BodyJson(map[string]interface{}{
    "key1": "value1",
  }).Do(ctx)
  return
}
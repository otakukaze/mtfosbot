package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"git.trj.tw/golang/mtfosbot/module/apis"
	"git.trj.tw/golang/mtfosbot/module/config"
)

type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ImageMessage struct {
	Type               string `json:"type"`
	OriginalContentUrl string `json:"originalContentUrl"`
	PreviewImageUrl    string `json:"previewImageUrl"`
}

type pushBody struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}
type replyBody struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

var baseUrl = "https://api.line.me/"

func getUrl(p string) (string, bool) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", false
	}
	ref, err := u.Parse(p)
	if err != nil {
		return "", false
	}
	str := ref.String()
	return str, true
}

func getHeaders() map[string]string {
	m := make(map[string]string)
	conf := config.GetConf()
	m["Content-Type"] = "application/json"
	m["Authorization"] = fmt.Sprintf("Bearer %s", conf.Line.Access)
	return m
}

// PushMessage -
func PushMessage(target string, message interface{}) {
	url := "/v2/bot/message/push"

	body := &pushBody{
		To: target,
	}

	switch message.(type) {
	case ImageMessage:
		break
	case TextMessage:
		break
	default:
		return
	}
	body.Messages = append(body.Messages, message)
	dataByte, err := json.Marshal(body)
	if err != nil {
		fmt.Println("json encoding error")
		return
	}

	byteReader := bytes.NewReader(dataByte)

	apiUrl, ok := getUrl(url)
	if !ok {
		fmt.Println("url parser fail")
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		Url:     apiUrl,
		Headers: getHeaders(),
		Body:    byteReader,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("post api fail")
		return
	}
}

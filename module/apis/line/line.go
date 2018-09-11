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

// TextMessage - line text message object
type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImageMessage - line image message object
type ImageMessage struct {
	Type               string `json:"type"`
	OriginalContentURL string `json:"originalContentUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
}

type pushBody struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}
type replyBody struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

var baseURL = "https://api.line.me/"

func getURL(p string) (string, bool) {
	u, err := url.Parse(baseURL)
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
	if len(target) == 0 {
		return
	}
	url := "/v2/bot/message/push"

	body := &pushBody{
		To: target,
	}

	switch message.(type) {
	case ImageMessage:
		m := (message.(ImageMessage))
		m.Type = "image"
		message = m
		break
	case TextMessage:
		m := (message.(TextMessage))
		m.Type = "text"
		message = m
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

	apiURL, ok := getURL(url)
	if !ok {
		fmt.Println("url parser fail")
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		Url:     apiURL,
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

// ReplyMessage -
func ReplyMessage(replyToken string, message interface{}) {
	if len(replyToken) == 0 {
		return
	}
	url := "/v2/bot/message/reply"

	body := &replyBody{
		ReplyToken: replyToken,
	}

	switch message.(type) {
	case ImageMessage:
		m := (message.(ImageMessage))
		m.Type = "image"
		message = m
		break
	case TextMessage:
		m := (message.(TextMessage))
		m.Type = "text"
		message = m
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

	apiURL, ok := getURL(url)
	if !ok {
		fmt.Println("url parser fail")
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		Url:     apiURL,
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

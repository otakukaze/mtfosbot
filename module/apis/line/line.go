package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

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

// LineUserInfo -
type LineUserInfo struct {
	DisplayName string `json:"displayName"`
	UserID      string `json:"userId"`
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
	log.Println("push target :::: ", target)
	if len(target) == 0 {
		return
	}
	urlPath := "/v2/bot/message/push"

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
		log.Println("to json error ::::", err)
		return
	}

	byteReader := bytes.NewReader(dataByte)

	apiURL, ok := getURL(urlPath)
	if !ok {
		log.Println("get url fail ::::::")
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		URL:     apiURL,
		Headers: getHeaders(),
		Body:    byteReader,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		log.Println("get req fail :::::: ", err)
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("send to line error :::: ", err)
		return
	}
}

// ReplyMessage -
func ReplyMessage(replyToken string, message interface{}) {
	if len(replyToken) == 0 {
		return
	}
	urlPath := "/v2/bot/message/reply"

	body := &replyBody{
		ReplyToken: replyToken,
	}

	switch message.(type) {
	case *ImageMessage:
		m := (message.(*ImageMessage))
		m.Type = "image"
		message = m
		break
	case *TextMessage:
		m := (message.(*TextMessage))
		m.Type = "text"
		message = m
		break
	default:
		fmt.Println("input type error")
		return
	}

	body.Messages = append(body.Messages, message)
	dataByte, err := json.Marshal(body)
	if err != nil {
		return
	}

	byteReader := bytes.NewReader(dataByte)

	apiURL, ok := getURL(urlPath)
	if !ok {
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		URL:     apiURL,
		Headers: getHeaders(),
		Body:    byteReader,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
}

// GetUserInfo -
func GetUserInfo(u, g string) (user *LineUserInfo, err error) {
	urlPath := fmt.Sprintf("/v2/bot/group/%s/member/%s", g, u)
	header := getHeaders()
	apiURL, ok := getURL(urlPath)
	if !ok {
		return nil, errors.New("url parser fail")
	}

	reqObj := apis.RequestObj{
		Method:  "GET",
		URL:     apiURL,
		Headers: header,
	}
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("api response not 200")
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, errors.New("response body not json")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return nil, err
	}

	return
}

// GetContentHead -
func GetContentHead(id string) (mime string, err error) {
	urlPath := fmt.Sprintf("/v2/bot/message/%s/content", id)
	header := getHeaders()
	u, ok := getURL(urlPath)
	if !ok {
		return "", errors.New("get url fail")
	}

	reqObj := apis.RequestObj{
		Method:  "HEAD",
		URL:     u,
		Headers: header,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	mime = resp.Header.Get("Content-Type")

	return
}

// DownloadContent -
func DownloadContent(id string, w io.Writer) (err error) {
	urlPath := fmt.Sprintf("/v2/bot/message/%s/content", id)
	header := getHeaders()
	u, ok := getURL(urlPath)
	if !ok {
		return errors.New("get url fail")
	}

	reqObj := apis.RequestObj{
		Method:  "GET",
		URL:     u,
		Headers: header,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)

	return
}

package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/line-message"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

// GetRawBody - line webhook body get
func GetRawBody(c *context.Context) {
	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.DataFormat("body read fail")
		return
	}
	c.Set("rawbody", byteBody)
	c.Next()
}

// VerifyLine - middleware
func VerifyLine(c *context.Context) {
	rawbody, ok := c.Get("rawbody")
	if !ok {
		c.DataFormat("body read fail")
		return
	}
	var raw []byte
	if raw, ok = rawbody.([]byte); !ok {
		c.DataFormat("body type error")
		return
	}
	sign := c.GetHeader("X-Line-Signature")
	if len(sign) == 0 {
		c.Next()
		return
	}

	conf := config.GetConf()

	hash := hmac.New(sha256.New, []byte(conf.Line.Secret))
	_, err := hash.Write(raw)
	if err != nil {
		c.ServerError(nil)
		return
	}
	hashSign := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if hashSign != sign {
		c.CustomRes(403, map[string]string{
			"message": "sign verify fail",
		})
		return
	}
	c.Next()
}

// GetLineMessage -
func GetLineMessage(c *context.Context) {
	rawbody, ok := c.Get("rawbody")
	if !ok {
		c.DataFormat("body read fail")
	}
	var raw []byte
	if raw, ok = rawbody.([]byte); !ok {
		c.DataFormat("body type error")
	}

	fmt.Println("Line Hook ::: ", string(raw))

	events := struct {
		Events []*lineobj.EventObject `json:"events"`
	}{}

	err := json.Unmarshal(raw, &events)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if len(events.Events) > 0 {
		for _, v := range events.Events {
			go linemsg.MessageEvent(v)
		}
	}

	c.Success(nil)
}

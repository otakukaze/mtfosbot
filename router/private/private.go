package private

import (
	"fmt"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
)

// VerifyKey -
func VerifyKey(c *context.Context) {
	conf := config.GetConf()
	key := c.GetHeader("X-Mtfos-Key")

	if len(key) == 0 {
		c.Forbidden(nil)
		return
	}

	if key != conf.SelfKey {
		c.Forbidden(nil)
		return
	}

	c.Next()
}

// GetFacebookPageIDs -
func GetFacebookPageIDs(c *context.Context) {
	pages, err := model.GetAllFacebookPage()
	if err != nil {
		c.ServerError(nil)
		return
	}

	ids := make([]string, 0)
	for _, v := range pages {
		ids = append(ids, v.ID)
	}

	c.Success(map[string]interface{}{
		"list": ids,
	})
}

// UpdateFacebookPagePost -
func UpdateFacebookPagePost(c *context.Context) {
	var err error
	type pageStruct struct {
		ID     string `json:"id"`
		PostID string `json:"post_id"`
		Link   string `json:"link"`
		Text   string `json:"text"`
	}
	bodyArg := struct {
		Pages []pageStruct `json:"pages"`
	}{}

	err = c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	for _, v := range bodyArg.Pages {
		if len(v.ID) == 0 || len(v.PostID) == 0 || len(v.Link) == 0 {
			continue
		}

		page, err := model.GetFacebookPage(v.ID)
		if err != nil {
			continue
		}
		if page.LastPost == v.PostID {
			continue
		}
		err = page.UpdatePost(v.PostID)
		if err != nil {
			continue
		}

		err = page.GetGroups()
		if err != nil {
			continue
		}

		for _, g := range page.Groups {
			if g.Notify {
				tmpl := g.Tmpl
				if len(tmpl) > 0 {
					tmpl = strings.Replace(tmpl, "{link}", v.Link, -1)
					tmpl = strings.Replace(tmpl, "{txt}", v.Text, -1)
				} else {
					tmpl = fmt.Sprintf("%s\n%s", v.Text, v.Link)
				}
				msg := line.TextMessage{
					Text: tmpl,
				}
				line.PushMessage(g.ID, msg)
			}
		}
	}

	c.Success(nil)
}

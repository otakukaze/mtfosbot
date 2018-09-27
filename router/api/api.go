package api

import (
	"strconv"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
	"github.com/gin-gonic/contrib/sessions"
	"golang.org/x/crypto/bcrypt"
)

// CheckSession -
func CheckSession(c *context.Context) {
	session := sessions.Default(c.Context)
	userData := session.Get("user")
	loginType := session.Get("loginType")
	if userData == nil || loginType == nil {
		c.LoginFirst(nil)
		return
	}
	var name string
	var ltype string
	var ok bool
	switch userData.(type) {
	case model.Account:
		name = userData.(model.Account).Account
	case twitch.UserInfo:
		name = userData.(twitch.UserInfo).Login
	default:
		c.LoginFirst(nil)
		return
	}
	if ltype, ok = loginType.(string); !ok {
		c.LoginFirst(nil)
		return
	}

	loginUser := map[string]string{
		"name": name,
		"type": ltype,
	}
	session.Set("loginUser", loginUser)
	session.Save()
	c.Next()
}

// UserLogin - system user login
func UserLogin(c *context.Context) {
	bodyArg := struct {
		Account  string `form:"account" json:"account" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}
	err := c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	acc, err := model.GetAccount(bodyArg.Account)
	if err != nil {
		c.ServerError(`account or password error`)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(bodyArg.Password))
	if err != nil {
		c.DataFormat(`account or password error`)
		return
	}

	session := sessions.Default(c.Context)

	acc.Password = ""
	session.Set("user", acc)
	session.Set("loginType", "system")
	session.Save()

	c.Success(nil)
}

// UserLogout -
func UserLogout(c *context.Context) {
	session := sessions.Default(c.Context)
	session.Clear()
	session.Save()

	c.Success(nil)
}

// GetSessionData -
func GetSessionData(c *context.Context) {
	session := sessions.Default(c.Context)
	loginUser := session.Get("loginUser")
	if loginUser == nil {
		c.LoginFirst(nil)
		return
	}

	user := map[string]interface{}{
		"user": loginUser,
	}
	c.Success(user)
}

// GetLineMessageLog -
func GetLineMessageLog(c *context.Context) {
	numP := 1
	if p, ok := c.GetQuery("p"); ok {
		if i, err := strconv.Atoi(p); err == nil {
			numP = i
		}
	}
	numMax := 20
	if max, ok := c.GetQuery("max"); ok {
		if m, err := strconv.Atoi(max); err == nil {
			numMax = m
		}
	}

	g := c.DefaultQuery("group", "")
	u := c.DefaultQuery("user", "")

	count, err := model.GetLineMessageLogCount()
	if err != nil {
		c.ServerError(nil)
		return
	}

	page := utils.CalcPage(count, numP, numMax)

	logs, err := model.GetLineMessageLogList(g, u, page.Offset, page.Limit)
	if err != nil {
		c.ServerError(nil)
		return
	}

	resMap := make([]map[string]interface{}, 0)

	for _, v := range logs {
		m := utils.ToMap(v.LineMessageLog)
		m["group_name"] = v.GroupName
		m["user_name"] = v.UserName
		resMap = append(resMap, m)
	}

	c.Success(map[string]interface{}{
		"list": resMap,
		"page": map[string]interface{}{
			"cur":   page.Page,
			"total": page.Total,
		},
	})
}

package api

import (
	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
	"github.com/gin-gonic/contrib/sessions"
	"golang.org/x/crypto/bcrypt"
)

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

	accInt := utils.ToMap(acc)
	delete(accInt, "password")
	session := sessions.Default(c.Context)

	session.Set("user", accInt)
	session.Set("loginType", "system")
	session.Save()

	c.Success(nil)
}

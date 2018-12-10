package api

import (
	"fmt"
	"strconv"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

// GetLineList -
func GetLineList(c *context.Context) {
	ls, err := model.GetLineGroupList()
	if err != nil {
		c.ServerError(nil)
		return
	}
	list := make([]map[string]interface{}, 0)

	if ls != nil {
		for _, v := range ls {
			if v == nil {
				continue
			}
			list = append(list, utils.ToMap(v))
		}
	}

	c.Success(map[string]interface{}{
		"list": list,
	})
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
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	where := make(map[string]string)
	if len(g) > 0 {
		where["group"] = g
	}
	if len(u) > 0 {
		where["user"] = u
	}

	count, err := model.GetLineMessageLogCount(where)
	if err != nil {
		fmt.Println("count error :::: ", err)
		c.ServerError(nil)
		return
	}

	page := utils.CalcPage(count, numP, numMax)
	if numP == -1 {
		numP = page.Total
	}
	page = utils.CalcPage(count, numP, numMax)

	logs, err := model.GetLineMessageLogList(g, u, page.Offset, page.Limit, order)
	if err != nil {
		c.ServerError(nil)
		return
	}

	resMap := make([]map[string]interface{}, len(logs))

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

// GetCommandList -
func GetCommandList(c *context.Context) {
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
	cmd := c.DefaultQuery("cmd", "")

	whereMap := make(map[string]string)

	if len(g) > 0 {
		whereMap["group"] = g
	}
	if len(cmd) > 0 {
		whereMap["cmd"] = cmd
	}

	count, err := model.GetCommandCount(whereMap)
	if err != nil {
		c.ServerError(nil)
		return
	}

	page := utils.CalcPage(count, numP, numMax)

	cmds, err := model.GetCommands(whereMap, page.Offset, page.Limit, nil)
	if err != nil {
		c.ServerError(nil)
		return
	}

	cmdList := make([]map[string]interface{}, len(cmds))
	for _, v := range cmds {
		tmp := utils.ToMap(v)
		tmp["group_name"] = v.GroupName
		cmdList = append(cmdList, tmp)
		tmp = nil
	}

	c.Success(map[string]interface{}{
		"list": cmdList,
		"page": map[string]interface{}{
			"cur":   page.Page,
			"total": page.Total,
		},
	})
}

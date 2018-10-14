package api

import (
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
		for v := range ls {
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

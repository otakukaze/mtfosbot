package msgcmd

import (
	"fmt"
	"regexp"
	"strings"

	"git.trj.tw/golang/mtfosbot/module/utils"

	"git.trj.tw/golang/mtfosbot/model"

	"git.trj.tw/golang/mtfosbot/module/apis/line"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

func parseCMD(in string) (c [][]string) {
	re, err := regexp.Compile("{{(.+?)}}")
	if err != nil {
		return
	}

	c = re.FindAllStringSubmatch(in, -1)

	return
}

// ParseLineMsg -
func ParseLineMsg(txt, replyToken string, source *lineobj.SourceObject) {
	fmt.Println("parse Line msg ;::: ", txt, utils.ToMap(source))
	if source.Type != "group" {
		return
	}
	strs := strings.Split(strings.Trim(txt, " "), " ")
	fmt.Println(strs)
	// skip empty string
	if len(strs[0]) == 0 {
		return
	}

	if strings.HasPrefix(strs[0], "!") || strings.HasPrefix(strs[0], "！") {
		// nor cmd
		cmd := strs[0][1:]
		fmt.Println(cmd)
		if len(cmd) == 0 {
			return
		}
		c, err := model.GetGroupCommand(strings.ToLower(cmd), source.GroupID)
		if err != nil || c == nil {
			fmt.Println("get command fail  ", err)
			return
		}
		fmt.Println("get command ::::", utils.ToMap(c))

		str := runCMD(strings.Join(strs[1:], " "), c.Message, source)
		fmt.Println("cmd result ::::::", str)
		m := parseResult(str)
		fmt.Println("parsed cmd result :::::", m)
		line.ReplyMessage(replyToken, m)

	} else {
		// key cmd
		c, err := model.GetGroupKeyCommand(strings.ToLower(strs[0]), source.GroupID)
		if err != nil || c == nil {
			fmt.Println("get command fail::::", err)
			return
		}
		fmt.Println("get command ::::", utils.ToMap(c))

		str := runCMD(strings.Join(strs[1:], " "), c.Message, source)
		fmt.Println("cmd result ::::", str)
		m := parseResult(str)
		fmt.Println("parsed cmd result :::::", m)
		line.ReplyMessage(replyToken, m)

	}
}

func parseResult(str string) interface{} {
	var m interface{}

	if strings.HasPrefix(str, "$image$") {
		str = strings.Replace(str, "$image$", "", 1)
		strs := strings.Split(str, ";")
		m = &line.ImageMessage{
			OriginalContentURL: strs[0],
			PreviewImageURL:    strs[1],
		}
	} else {
		m = &line.TextMessage{
			Text: str,
		}
	}

	return m
}

func runCMD(txt, c string, s *lineobj.SourceObject) (res string) {
	fmt.Println("run run run::::::", txt, c, s)
	cmdAct := parseCMD(c)
	if len(cmdAct) == 0 {
		return c
	}
	res = c
	for _, v := range cmdAct {
		if len(v) > 1 {
			// run cmd
			m := strings.Split(v[1], "=")
			sub := ""
			if len(m) > 1 {
				sub = strings.Join(m[1:], " ")
			}
			fmt.Println("run ::::::::", m[0], sub, txt)
			cmdRes := selectAct(m[0], sub, txt, s)
			res = strings.Replace(res, v[0], cmdRes, 1)
		}
	}
	return
}
package linemsg

import (
	"fmt"
	"os"
	"path"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/config"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
	msgcmd "git.trj.tw/golang/mtfosbot/module/message-command"
)

func messageType(e *lineobj.EventObject) {
	msg := e.Message
	mtype, ok := msg["type"]
	if !ok {
		return
	}

	if t, ok := mtype.(string); ok {
		switch t {
		case "text":
			textMsg(e)
			break
		case "image":
			imageMsg(e)
			break
		}
	}
	return
}

func textMsg(e *lineobj.EventObject) {
	msg := e.Message
	mtxt, ok := msg["text"]
	if !ok {
		return
	}

	// group action
	if e.Source.Type == "group" {
		if txt, ok := mtxt.(string); ok {
			msgcmd.ParseLineMsg(txt, e.ReplyToken, e.Source)
			saveTextMsgToLog(txt, e.Source)
		}
	}
	return
}

func imageMsg(e *lineobj.EventObject) {
	msg := e.Message
	imgID, ok := msg["id"]
	if !ok {
		return
	}
	// group action
	if e.Source.Type == "group" {
		if id, ok := imgID.(string); ok {
			saveImageMsgToLog(id, e.Source)
		}
	}
}

func getSourceUser(uid, gid string) (u *model.LineUser, err error) {
	userData, err := model.GetLineUserByID(uid)
	if err != nil {
		return
	}
	if userData == nil {
		tmpu, err := line.GetUserInfo(uid, gid)
		if err != nil || tmpu == nil {
			return nil, err
		}
		userData = &model.LineUser{}
		userData.ID = tmpu.UserID
		userData.Name = tmpu.DisplayName
		err = userData.Add()
		if err != nil {
			return nil, err
		}
	}

	return userData, nil
}

func saveTextMsgToLog(txt string, s *lineobj.SourceObject) {
	u, err := getSourceUser(s.UserID, s.GroupID)
	if err != nil || u == nil {
		return
	}

	model.AddLineMessageLog(s.GroupID, s.UserID, txt, "text")
}

func saveImageMsgToLog(id string, s *lineobj.SourceObject) {
	u, err := getSourceUser(s.UserID, s.GroupID)
	if err != nil || u == nil {
		return
	}

	mime, err := line.GetContentHead(id)
	if err != nil || len(mime) == 0 {
		return
	}

	ext := ""
	switch mime {
	case "image/jpeg":
		ext = ".jpg"
		break
	case "image/jpg":
		ext = ".jpg"
		break
	case "image/png":
		ext = ".png"
		break
	default:
		return
	}

	conf := config.GetConf()

	fname := fmt.Sprintf("log_%s%s", id, ext)

	fullPath := path.Join(conf.LogImageRoot, fname)

	w, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer w.Close()

	err = line.DownloadContent(id, w)
	if err != nil {
		return
	}

	model.AddLineMessageLog(s.GroupID, s.UserID, fname, "image")
}

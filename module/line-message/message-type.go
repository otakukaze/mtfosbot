package linemsg

import (
	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
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

func saveTextMsgToLog(txt string, s *lineobj.SourceObject) {
	userData, err := model.GetLineUserByID(s.UserID)
	if err != nil {
		return
	}
	if userData == nil {
		tmpu, err := line.GetUserInfo(s.UserID, s.GroupID)
		if err != nil || tmpu == nil {
			return
		}
		userData = &model.LineUser{}
		userData.ID = tmpu.UserID
		userData.Name = tmpu.DisplayName
		err = userData.Add()
		if err != nil {
			return
		}
	}

	model.AddLineMessageLog(s.GroupID, s.UserID, txt)
}

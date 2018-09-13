package linemsg

import (
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
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
		if _, ok := mtxt.(string); ok {

		}
	}
	return
}

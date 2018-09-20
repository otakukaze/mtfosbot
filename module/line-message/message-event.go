package linemsg

import (
	"fmt"

	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

// MessageEvent -
func MessageEvent(e *lineobj.EventObject) {

	switch e.Type {
	case "message":
		messageType(e)
		break
	default:
		fmt.Println("line webhook type not match")
	}
}

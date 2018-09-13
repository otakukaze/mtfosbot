package linemsg

import (
	"fmt"

	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

// MessageEvent -
func MessageEvent(e *lineobj.EventObject) {
	fmt.Println(utils.ToMap(e))

	switch e.Type {
	case "message":
		messageType(e)
		break
	}
}

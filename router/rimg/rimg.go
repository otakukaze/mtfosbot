package rimg

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/nfnt/resize"

	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

// GetOriginImage -
func GetOriginImage(c *context.Context) {
	fname := c.Param("imgname")
	conf := config.GetConf()
	imgP := utils.ParsePath(conf.ImageRoot)
	exists := utils.CheckExists(imgP, true)
	if !exists {
		c.NotFound("image path not found")
		return
	}

	newP := path.Join(imgP, fname)
	exists = utils.CheckExists(newP, false)
	if !exists {
		c.NotFound("image file not found")
		return
	}

	fileBuf, err := os.Open(newP)
	if err != nil {
		c.ServerError(nil)
		return
	}
	defer fileBuf.Close()

	imgf, _, err := image.DecodeConfig(fileBuf)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if imgf.Height > 1024 || imgf.Width > 1024 {
		img, _, err := image.Decode(fileBuf)
		if err != nil {
			c.ServerError(nil)
			return
		}
		w := 0
		h := 0
		if imgf.Width > imgf.Height {
			w = 1024
		} else {
			h = 1024
		}
		buf := new(bytes.Buffer)
		m := resize.Resize(uint(w), uint(h), img, resize.Bilinear)
		err = jpeg.Encode(buf, m, nil)
		if err != nil {
			c.ServerError(nil)
			return
		}
		c.Writer.Header().Set("Content-Type", "image/jpeg")
		c.Writer.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		_, err = c.Writer.Write(buf.Bytes())
		if err != nil {
			c.ServerError(nil)
		}
		return
	}

	buf, err := ioutil.ReadAll(fileBuf)
	c.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	_, err = c.Writer.Write(buf)
	if err != nil {
		c.ServerError(nil)
		return
	}
}

package utils

import (
	"math"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// PageObject -
type PageObject struct {
	Page   int `json:"page" cc:"page"`
	Total  int `json:"total" cc:"total"`
	Offset int `json:"offset" cc:"offset"`
	Limit  int `json:"limit" cc:"limit"`
}

// CalcPage -
func CalcPage(count, page, max int) (po PageObject) {
	if count < 0 {
		count = 0
	}
	if page < 1 {
		page = 1
	}
	if max < 1 {
		max = 1
	}

	total := int(math.Ceil(float64(count) / float64(max)))
	if total < 1 {
		total = 1
	}
	if page > total {
		page = total
	}
	offset := (page - 1) * max
	if offset > count {
		offset = count
	}
	limit := max

	po = PageObject{}
	po.Limit = limit
	po.Page = page
	po.Offset = offset
	po.Total = total

	return
}

// ToMap struct to map[string]interface{}
func ToMap(ss interface{}) map[string]interface{} {
	t := reflect.ValueOf(ss)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	smap := make(map[string]interface{})
	mtag := regexp.MustCompile(`cc:\"(.+)\"`)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := string(t.Type().Field(i).Tag)
		str := mtag.FindStringSubmatch(tag)
		name := t.Type().Field(i).Name
		if len(str) > 1 {
			name = str[1]
		}
		if name != "-" {
			smap[name] = f.Interface()
		}
	}

	return smap
}

// ParsePath - parse file path to absPath
func ParsePath(dst string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	if []rune(dst)[0] == '~' {
		home := UserHomeDir()
		if len(home) > 0 {
			dst = strings.Replace(dst, "~", home, -1)
		}
	}

	if path.IsAbs(dst) {
		dst = path.Clean(dst)
		return dst
	}

	str := path.Join(wd, dst)
	str = path.Clean(str)
	return str
}

// UserHomeDir - get user home directory
func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// CheckExists - check file exists
func CheckExists(filePath string, allowDir bool) bool {
	filePath = ParsePath(filePath)
	stat, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return false
	}
	if !allowDir && stat.IsDir() {
		return false
	}
	return true
}

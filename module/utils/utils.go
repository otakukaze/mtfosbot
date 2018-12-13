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
		leveling := false
		if len(str) > 1 {
			strArr := strings.Split(str[1], ",")
			name = strArr[0]
			if len(strArr) > 1 && strArr[1] == "<<" {
				leveling = true
			}
		}

		if f.Kind() == reflect.Slice {
			if name == "-" {
				continue
			}
			if f.Len() > 0 && f.Index(0).Kind() != reflect.Struct {
				smap[name] = f.Interface()
				continue
			}

			tmp := make([]map[string]interface{}, f.Len())
			for i := 0; i < f.Len(); i++ {
				tmp[i] = ToMap(f.Index(i).Interface())
			}
			smap[name] = tmp
			continue
		}

		if f.Kind() == reflect.Struct {
			if name == "-" && leveling == false {
				continue
			}

			if leveling == true {
				tmp := ToMap(f.Interface())
				for k, v := range tmp {
					smap[k] = v
				}
			} else if name != "-" && leveling == false {
				smap[name] = f.Interface()
			}
			continue
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

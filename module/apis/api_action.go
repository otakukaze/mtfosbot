package apis

import (
	"io"
	"net/http"
)

type RequestObj struct {
	Method  string
	Url     string
	Body    io.Reader
	Headers map[string]string
}

func GetRequest(r RequestObj) (req *http.Request, err error) {

	req, err = http.NewRequest(r.Method, r.Url, r.Body)
	if err != nil {
		return
	}

	if len(r.Headers) > 0 {
		for k, v := range r.Headers {
			req.Header.Set(k, v)
		}
	}

	return
}

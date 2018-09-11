package twitch

import (
	"fmt"
	"net/url"
)

var baseURL = "https://api.twitch.tv"

func getURL(p string) (string, bool) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}
	ref, err := u.Parse(p)
	if err != nil {
		return "", false
	}
	str := ref.String()
	return str, true
}

func getHeaders(token string) map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"
	m["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return m
}

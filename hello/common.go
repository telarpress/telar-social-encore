package hello

import (
	"net/http"
	"strings"
)

// RemoveBaseURLFromRequest removes the base url from the request path
func RemoveBaseURLFromRequest(r *http.Request) {
	r.URL.Path = RemoveBaseURL(r.URL.Path)
	r.RequestURI = RemoveBaseURL(r.RequestURI)
}

// RemoveBaseURL removes the base url from the request path
func RemoveBaseURL(url string) string {
	if "/" == url {
		return url
	}
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}
	return "/"+strings.Join(strings.Split(url, "/")[1:], "/")
}

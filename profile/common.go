package profile

import (
	"net/http"
	"strings"
)

// RemoveBaseURLFromRequest removes the base url from the request path
func RemoveBaseURLFromRequest(r *http.Request) {
	r.URL.Path = RemoveBaseURL(r.URL.Path)
	r.RequestURI = RemoveBaseURL(r.RequestURI)
	// If Encore is running in local development the remote address is not a valid hostname,
	// which Fiber expects. Set this to a valid placeholder value to appease Fiber.
	if strings.Contains(r.RemoteAddr, "yamux") {
		r.RemoteAddr = "localhost:0"
	}
}

// RemoveBaseURL removes the base url from the request path
func RemoveBaseURL(url string) string {
	if url == "/" {
		return url
	}
	url = strings.TrimPrefix(url, "/")
	return "/" + strings.Join(strings.Split(url, "/")[1:], "/")
}

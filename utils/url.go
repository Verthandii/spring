package utils

import (
	"fmt"
	"net/http"
)

func Url(request *http.Request) string {
	if request == nil {
		return ""
	}

	url := "https://"

	if request.TLS == nil {
		url = "http://"
	}

	return fmt.Sprintf("%s%s%s", url, request.Host, request.URL.String())
}

package usermanagerclient

import (
	"io"
	"io/ioutil"
)

func getBodyString(body io.ReadCloser) string {
	bodyBytes, _ := ioutil.ReadAll(body)
	return string(bodyBytes)
}

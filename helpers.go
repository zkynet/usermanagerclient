package usermanagerclient

import (
	"io"
	"io/ioutil"
	"time"
)

func getBodyString(body io.ReadCloser) string {
	bodyBytes, _ := ioutil.ReadAll(body)
	return string(bodyBytes)
}

func (c *System) isCookieExpired() bool {

	for _, v := range c.Cookies {
		if v.Name == c.SystemCookieName {
			return v.Expires.Before(time.Now().Add(time.Duration(10) * time.Minute))
		}
	}

	return false
}

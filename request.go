package usermanagerclient

import (
	"bytes"
	"net/http"
)

func (c *System) requestWithSystemCredentials(headers map[string]string, method string, payload []byte, domain string) (error, *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, domain, bytes.NewReader(payload))
	if err != nil {
		return err, nil
	}

	for i, v := range headers {
		req.Header.Add(i, v)
	}

	for _, v := range c.Cookies {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	for _, v := range resp.Cookies() {
		c.Cookies[v.Name] = v
	}

	for i, v := range resp.Header {
		if i == c.JWTHeaderName {
			c.Headers[c.JWTHeaderName] = v[0]
		}
	}

	return nil, resp
}

func (s *System) requestWithUserCredentials(headers map[string]string, method string, payload []byte, domain string, cookie *http.Cookie, jwt string) (error, *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, domain, bytes.NewReader(payload))
	if err != nil {
		return err, nil
	}

	if jwt != "" {
		req.Header.Add(s.JWTHeaderName, jwt)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	return nil, resp
}

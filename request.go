package usermanagerclient

import (
	"bytes"
	"net/http"
)

func (c *Client) sendRequest(headers map[string]string, method string, payload []byte, domain string) (error, *http.Response) {
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

	return nil, resp
}

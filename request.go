package usermanagerclient

import (
	"bytes"
	"net/http"
)

func sendRequest(headers map[string]string, method string, payload []byte, domain string) (error, *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, domain, bytes.NewReader(payload))
	if err != nil {
		return err, nil
	}

	for i, v := range headers {
		req.Header.Add(i, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	return nil, resp
}

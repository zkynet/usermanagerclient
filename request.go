package usermanagerclient

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
)

func (c *System) requestWithSystemCredentials(headers map[string]string, method string, payload []byte, domain string) (error, *http.Response) {
	log.Println("sending to:", domain)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	req, err := http.NewRequest(method, domain, bytes.NewReader(payload))
	if err != nil {
		return err, nil
	}

	for i, v := range headers {
		req.Header.Add(i, v)
	}

	req.Header.Add(c.SystemAuthHeaderKey, c.SystemAuthHeaderValue)

	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	for i, v := range resp.Header {
		if i == c.JWTHeaderKey {
			c.Headers[c.JWTHeaderKey] = v[0]
		}
	}

	return nil, resp
}

func (s *System) requestWithUserCredentials(headers map[string]string, method string, payload []byte, domain string, jwt string) (error, *http.Response) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	req, err := http.NewRequest(method, domain, bytes.NewReader(payload))
	if err != nil {
		return err, nil
	}

	req.Header.Add(s.SystemAuthHeaderKey, s.SystemAuthHeaderValue)

	if jwt != "" {
		req.Header.Add(s.JWTHeaderKey, jwt)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	return nil, resp
}

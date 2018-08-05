package usermanagerclient

import (
	"encoding/json"
	"log"
)

type Client struct {
	URL     string
	Port    string
	Headers map[string]string
}

func (c *Client) Create(name string, phone string, email string, password string) error {
	message := map[string]interface{}{
		"name":     name,
		"phone":    phone,
		"email":    email,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	sendRequest(c.Headers, "POST", bytesRepresentation, c.URL+c.Port)
}

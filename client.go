package usermanagerclient

import (
	"encoding/json"
	"fmt"
	"log"
)

type Client struct {
	URL     string
	Port    string
	Headers map[string]string
}

func NewClient() *Client {
	return &Client{
		Headers: make(map[string]string),
	}
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

	url := c.URL + ":" + c.Port + "/user"
	err, resp := sendRequest(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return err
}

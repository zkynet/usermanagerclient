package usermanagerclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	URL     string
	Port    string
	Headers map[string]string
	Cookies map[string]*http.Cookie
}

func NewClient() *Client {
	return &Client{
		Headers: make(map[string]string),
		Cookies: make(map[string]*http.Cookie),
	}
}

func (c *Client) FacebookLogin(email string, name string, facebookID string) error {
	message := map[string]interface{}{
		"facebook_id": facebookID,
		"email":       email,
		"name":        name,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	url := c.URL + ":" + c.Port + "/facebook/login"
	err, resp := c.sendRequest(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return err
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
	err, resp := c.sendRequest(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return err
}

func (c *Client) ValidateRequest(namespace string, token string) error {
	message := map[string]interface{}{
		"tag": namespace,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	url := c.URL + ":" + c.Port + "/user"
	err, resp := c.sendRequest(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return err

}

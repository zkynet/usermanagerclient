package usermanagerclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type System struct {
	URL                   string
	Port                  string
	Headers               map[string]string
	Cookies               map[string]*http.Cookie
	JWT                   string
	SystemAuthHeaderKey   string
	SystemAuthHeaderValue string
	JWTHeaderKey          string
	SystemCookieName      string
	UserCookieName        string
}

type User struct {
	System *System
	JWT    string
	Cookie *http.Cookie
}

type Group struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type UserID struct {
	ID string `json:"id"`
}

func (c *System) Logout() error {

	url := c.URL + ":" + c.Port + "/logout"
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", nil, url)
	if err != nil {
		return err
	}

	fmt.Println(getBodyString(resp.Body))
	return err
}

func (c *System) Login(email string, password string) (error, string, int) {

	message := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/login"
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err, "", 0
	}

	data := &UserID{}
	json.NewDecoder(resp.Body).Decode(data)
	return err, data.ID, resp.StatusCode
}

func (c *System) CreateUser(name string, phone string, email string, password string, facebookID string, appID string) (error, string, int) {
	message := map[string]interface{}{
		"name":        name,
		"phone":       phone,
		"email":       email,
		"password":    password,
		"facebook_id": facebookID,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/user/" + appID
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err, "", 0
	}

	data := &UserID{}
	json.NewDecoder(resp.Body).Decode(data)
	return err, data.ID, resp.StatusCode
}

func (c *System) ValidateRequest(namespace string, request *http.Request) (error, string, int) {
	message := map[string]interface{}{
		"tag": namespace,
	}

	var userCookie *http.Cookie
	for _, cookie := range request.Cookies() {
		if cookie.Name == c.UserCookieName {
			userCookie = cookie
		}
	}

	var userJWT string
	for headerIndex, v := range request.Header {
		if headerIndex == c.JWTHeaderKey {
			userJWT = v[0]
		}
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/validateRequest"
	err, resp := c.requestWithUserCredentials(c.Headers, "POST", bytesRepresentation, url, userCookie, userJWT)
	if err != nil {
		return err, "", 0
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return err, string(bodyBytes), resp.StatusCode
}

func (c *System) CreateGroup(tag string, appID string) (error, string, int) {

	message := map[string]interface{}{
		"tag": tag,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/group/" + appID
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err, "", 0
	}

	data := &Group{}
	json.NewDecoder(resp.Body).Decode(data)

	return err, data.ID, resp.StatusCode
}
func (c *System) AttachNamespaceToGroup(groupID string, namespaceID string) (error, string, int) {

	url := c.URL + ":" + c.Port + "/namespace/" + namespaceID + "/attach/to/group/" + groupID
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", nil, url)
	if err != nil {
		return err, "", 0
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return err, string(bodyBytes), resp.StatusCode
}
func (c *System) AttachUserToGroup(groupID string, userID string) (error, string, int) {

	url := c.URL + ":" + c.Port + "/user/" + userID + "/attach/to/group/" + groupID
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", nil, url)
	if err != nil {
		return err, "", 0
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return err, string(bodyBytes), resp.StatusCode
}
func (c *System) CreateNamespace(tag string, appID string) (error, string, int) {

	message := map[string]interface{}{
		"tag": tag,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/namespace/" + appID
	err, resp := c.requestWithSystemCredentials(c.Headers, "POST", bytesRepresentation, url)
	if err != nil {
		return err, "", 0
	}

	data := &Group{}
	json.NewDecoder(resp.Body).Decode(data)

	return err, data.ID, resp.StatusCode
}

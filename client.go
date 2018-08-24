package usermanagerclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zkynet/errorwrapper"
)

type System struct {
	URL                   string
	Port                  string
	Headers               map[string]string
	JWT                   string
	SystemAuthHeaderKey   string
	SystemAuthHeaderValue string
	JWTHeaderKey          string
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

func (c *System) Login(email string, password string, appID string) (error, string, int) {

	message := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err, "", 0
	}

	url := c.URL + ":" + c.Port + "/login/" + appID
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

func (c *System) ValidateRequest(namespace string, request *http.Request) error {
	return c.ParseUserManagerError(
		c.ValidateNamspace(namespace, request))
}

func (c *System) ValidateNamspace(namespace string, request *http.Request) (error, string, int) {
	message := map[string]interface{}{
		"tag": namespace,
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
	err, resp := c.requestWithUserCredentials(c.Headers, "POST", bytesRepresentation, url, userJWT)
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

func (c *System) ParseUserManagerError(err error, body string, code int) error {
	if code == 200 {
		return nil
	}

	if code == 401 {
		return errorwrapper.Unauthorized(err)
	}

	if code == 403 {
		return errorwrapper.UnauthorizedCustomMessage(err, "You do not have the correct authentication headers")
	}

	newErr := errorwrapper.GenericError(err)
	newErr.Message = body
	newErr.OriginalError = err
	return newErr

}

func (c *System) GINROUTERNamespaceValidation(namespace string) gin.HandlerFunc {
	return func(g *gin.Context) {
		authErr := c.ParseUserManagerError(c.ValidateNamspace(namespace, g.Request))
		if authErr != nil {
			g.JSON(errorwrapper.GetErrorCode(authErr), errorwrapper.HandleError(authErr))
			g.Abort()
			return
		}
		g.Next()
	}
}

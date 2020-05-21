package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	oauthClientID string
	oauthSecret   string
	breezeURL     string
)

type BreezeClient struct {
	BaseUrl         *url.URL
	*http.Client
}

var Api *BreezeClient

func init() {
	baseUrl, _ := url.Parse(breezeURL)
	Api = &BreezeClient{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		BaseUrl: baseUrl,
	}
}

type passwordGrantLogin struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Scope        string `json:"scope"`
}

func (c *BreezeClient) Login(email string, password string) (string, error) {
	reqData := &passwordGrantLogin{
		GrantType:    "password",
		ClientId:     oauthClientID,
		ClientSecret: oauthSecret,
		Username:     email,
		Password:     password,
		Scope:        "*",
	}

	b, err := json.Marshal(reqData)
	rawData := bytes.NewBuffer(b)
	reqUrl := c.BaseUrl.ResolveReference(&url.URL{Path: "oauth/token"})
	res, err := c.Post(reqUrl.String(), "application/json", rawData)

	if err != nil {
		return "", err
	}

	resData := map[string]string{}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("request failed with code %d", res.StatusCode)
	}
	b, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &resData)

	return resData["access_token"], nil
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/virtualops/cli/pkg/config"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	oauthSecret string
	breezeURL   string
)

type BreezeClient struct {
	unauthenticated bool
	*http.Client
}

func (c *BreezeClient) WithoutAuth(callback func(c *BreezeClient)) {
	shallowClientClone := *c
	shallowClientClone.unauthenticated = true
	callback(&shallowClientClone)
	shallowClientClone.unauthenticated = false
}

func (c *BreezeClient) Do(request *http.Request) (*http.Response, error) {
	if ! c.unauthenticated {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GlobalConfig.AuthToken))
	}

	//request.URL

	return c.Client.Do(request)
}

var api = &BreezeClient{
	Client: &http.Client{
		Timeout: time.Second * 10,
	},
}

func init() {
}

func Login(email string, password string) (string, error) {
	reqData := bytes.NewBufferString(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	res, err := api.Post("/api/login", "application/json", reqData)

	if err != nil {
		return "", err
	}

	resData := map[string]string{}
	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &resData)

	return resData["token"], nil
}

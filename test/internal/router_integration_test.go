// +build integration

package api_test

import (
	"chatter/internal/db"
	"chatter/model"
	"chatter/test"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	commonTest "github.com/EmotivesProject/common/test"
	"github.com/stretchr/testify/assert"
)

const uaclCreateUserEndpoint = "http://0.0.0.0:8082/user"

func TestRouterCreateToken(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, uaclCreateUserEndpoint)

	req, _ := http.NewRequest("GET", test.TS.URL+"/ws_token", nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterMessages(t *testing.T) {
	test.SetUpIntegrationTest()

	username, token := commonTest.CreateNewUser(t, uaclCreateUserEndpoint)

	username2, _ := commonTest.CreateNewUser(t, uaclCreateUserEndpoint)

	db.CreateUser(context.Background(), model.User{
		Username:  strings.ToLower(username),
		UserGroup: "test",
	})
	db.CreateUser(context.Background(), model.User{
		Username:  strings.ToLower(username2),
		UserGroup: "test",
	})

	url := fmt.Sprintf("%s/messages?to=%s&from=%s", test.TS.URL, strings.ToLower(username2), strings.ToLower(username))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterConnected(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, uaclCreateUserEndpoint)

	url := fmt.Sprintf("%s/connections", test.TS.URL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterWS(t *testing.T) {
	test.SetUpIntegrationTest()

	username, token := commonTest.CreateNewUser(t, uaclCreateUserEndpoint)

	url := fmt.Sprintf("%s/ws?token=%s", test.TS.URL, username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusForbidden)

	test.TearDownIntegrationTest()
}

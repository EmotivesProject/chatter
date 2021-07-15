// +build integration

package api_test

import (
	"chatter/internal/db"
	"chatter/test"
	"context"
	"fmt"
	"net/http"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterCreateToken(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	req, _ := http.NewRequest("GET", test.TS.URL+"/ws_token", nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterMessages(t *testing.T) {
	test.SetUpIntegrationTest()

	username, token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	db.CreateUser(context.Background(), username)

	url := fmt.Sprintf("%s/messages?to=%s&from=%s", test.TS.URL, "tom", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterConnected(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	url := fmt.Sprintf("%s/connections", test.TS.URL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterWS(t *testing.T) {
	test.SetUpIntegrationTest()

	username, token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	url := fmt.Sprintf("%s/ws?token=%s", test.TS.URL, username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusForbidden)

	test.TearDownIntegrationTest()
}

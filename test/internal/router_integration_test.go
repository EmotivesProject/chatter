// +build integration

package api_test

import (
	"chatter/test"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
)

func TestRouter(t *testing.T) {
	test.SetUpIntegrationTest()

	commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	test.TearDownIntegrationTest()
}

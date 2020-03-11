package testing

import (
	"fmt"
	"gcloud/gcorecloud-go/gcore"
	"gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/identity/tokens"
	"gcloud/gcorecloud-go/testhelper"
)

// authTokenPost verifies that providing certain AuthOptions results in an expected JSON structure.
func authTokenPost(t *testing.T, options gcorecloud.AuthOptions, requestJSON string) *tokens.Token {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	provider, err := gcore.NewGCoreClient(testhelper.Endpoint())
	require.NoError(t, err)
	serviceClient, err := gcore.NewIdentity(provider, gcorecloud.EndpointOpts{})
	require.NoError(t, err)

	testhelper.Mux.HandleFunc("/auth/jwt/login", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprint(w, TokenOutput)
		if err != nil {
			log.Error(err)
		}
	})

	actual, err := tokens.Create(serviceClient, &options).ExtractTokens()
	require.NoError(t, err)
	require.Equal(t, expectedToken, *actual)
	return actual
}

func authTokenPostErr(t *testing.T, options gcorecloud.AuthOptions, includeToken bool, expectedErr error) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	provider, err := gcore.NewGCoreClient(testhelper.Endpoint())
	require.NoError(t, err)
	serviceClient, err := gcore.NewIdentity(provider, gcorecloud.EndpointOpts{})
	require.NoError(t, err)

	if includeToken {
		serviceClient.AccessTokenID = "abcdef123456"
	}

	testhelper.Mux.HandleFunc("/auth/jwt/login", func(w http.ResponseWriter, r *http.Request) {

		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")

		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, `{"error": ""}`)
		if err != nil {
			log.Error(err)
		}
	})

	_, err = tokens.Create(serviceClient, &options).ExtractTokens()
	require.Error(t, err)
	require.IsType(t, err, expectedErr)
}

func TestCreateUserIDAndPassword(t *testing.T) {
	authTokenPost(t, gcorecloud.AuthOptions{Username: "me", Password: "squirrel!"}, `
		{
			"username": "me",
			"password": "squirrel!"
		}
	`)
}

func TestCreateExtractsTokenFromResponse(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	options := gcorecloud.AuthOptions{Username: "me", Password: "shhh"}
	token := authTokenPost(t, options, `
		{
			"username": "me",
			"password": "shhh"
		}
	`)

	require.Equal(t, token.Access, client.AccessToken)

}

func TestCreateFailureEmptyAuth(t *testing.T) {
	authTokenPostErr(t, gcorecloud.AuthOptions{}, false, gcorecloud.ErrDefault400{})
}

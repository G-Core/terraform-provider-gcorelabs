package testing

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

var testURL = "/v1/magnum/"

func TestAuthenticatedClient(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/auth/jwt/login", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "Authorization", "")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, client.AccessToken, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	options := gcorecloud.AuthOptions{
		Username: "me",
		Password: "secret",
		APIURL:   th.GCoreIdentifyEndpoint(),
		AuthURL:  th.GCoreRefreshTokenIdentifyEndpoint(),
	}
	provider, err := gcore.AuthenticatedClient(options)
	require.NoError(t, err)
	require.Equal(t, client.AccessToken, provider.AccessToken())
	require.Equal(t, client.RefreshToken, provider.RefreshToken())
}

func TestReauthAuthenticatedServiceClient(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	updatedAccessToken := client.AccessToken + "X"
	reauthCount := 0

	th.Mux.HandleFunc("/auth/jwt/login", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "Authorization", "")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, client.AccessToken, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	th.Mux.HandleFunc("/auth/jwt/refresh", func(w http.ResponseWriter, r *http.Request) {
		reauthCount++
		th.TestHeader(t, r, "Authorization", "")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, updatedAccessToken, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	serviceClient := client.ServiceAuthClient("magnum", "v1")
	fullTestURL := serviceClient.ResourceBaseURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), " ")[1]
		require.Contains(t, []string{client.AccessToken, updatedAccessToken}, header)
		if header == client.AccessToken {
			w.WriteHeader(http.StatusUnauthorized)
		} else if header == updatedAccessToken {
			w.WriteHeader(http.StatusOK)
		}
		_, err := fmt.Fprintf(w, `{}`)
		if err != nil {
			log.Error(err)
		}
	})

	require.Equal(t, client.AccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())
	r := gcorecloud.Result{}

	resp, err := serviceClient.Get(fullTestURL, &r.Body, nil)
	require.NoError(t, err)

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	require.Equal(t, resp.StatusCode, 200)

	require.Equal(t, updatedAccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())

	// retry
	serviceClient.AccessTokenID = client.AccessToken
	require.Equal(t, client.AccessToken, serviceClient.AccessToken())

	resp, err = serviceClient.Get(fullTestURL, &r.Body, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)

	require.Equal(t, updatedAccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())

	require.Equal(t, reauthCount, 2)

}

func TestReauthAuthenticatedServiceClientWithBadRefreshToken(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	updatedAccessToken := client.AccessToken + "X"
	reauthCount := 0
	authCount := 0

	th.Mux.HandleFunc("/auth/jwt/login", func(w http.ResponseWriter, r *http.Request) {
		authCount++
		th.TestHeader(t, r, "Authorization", "")
		token := client.AccessToken
		if authCount > 1 {
			token = updatedAccessToken
		}
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, token, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	th.Mux.HandleFunc("/auth/jwt/refresh", func(w http.ResponseWriter, r *http.Request) {
		reauthCount++
		th.TestHeader(t, r, "Authorization", "")
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, updatedAccessToken, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	serviceClient := client.ServiceAuthClient("magnum", "v1")
	fullTestURL := serviceClient.ResourceBaseURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), " ")[1]
		require.Contains(t, []string{client.AccessToken, updatedAccessToken}, header)
		if header == client.AccessToken {
			w.WriteHeader(http.StatusUnauthorized)
		} else if header == updatedAccessToken {
			w.WriteHeader(http.StatusOK)
		}
		_, err := fmt.Fprintf(w, `{}`)
		if err != nil {
			log.Error(err)
		}
	})

	require.Equal(t, client.AccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())
	r := gcorecloud.Result{}

	resp, err := serviceClient.Get(fullTestURL, &r.Body, nil) // nolint
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)

	require.Equal(t, updatedAccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())

	require.Equal(t, reauthCount, 1)
	require.Equal(t, authCount, 2)

}

func TestReauthTokenServiceClient(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	updatedAccessToken := client.AccessToken + "X"
	reauthCount := 0

	th.Mux.HandleFunc("/v1/token/refresh", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "Authorization", "")
		reauthCount++
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{ "access": "%s", "refresh": "%s"}`, updatedAccessToken, client.RefreshToken)
		if err != nil {
			log.Error(err)
		}
	})

	serviceClient := client.ServiceTokenClient("magnum", "v1")
	fullTestURL := serviceClient.ResourceBaseURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), " ")[1]
		require.Contains(t, []string{client.AccessToken, updatedAccessToken}, header)
		if header == client.AccessToken {
			w.WriteHeader(http.StatusUnauthorized)
		} else if header == updatedAccessToken {
			w.WriteHeader(http.StatusOK)
		}
		_, err := fmt.Fprintf(w, `{}`)
		if err != nil {
			log.Error(err)
		}
	})

	require.Equal(t, client.AccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())
	r := gcorecloud.Result{}

	resp, err := serviceClient.Get(fullTestURL, &r.Body, nil) // nolint
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)

	require.Equal(t, updatedAccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())

	// retry
	serviceClient.AccessTokenID = client.AccessToken
	require.Equal(t, client.AccessToken, serviceClient.AccessToken())

	resp, err = serviceClient.Get(fullTestURL, &r.Body, nil) // nolint
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)

	require.Equal(t, updatedAccessToken, serviceClient.AccessToken())
	require.Equal(t, client.RefreshToken, serviceClient.RefreshToken())

	require.Equal(t, reauthCount, 2)

}

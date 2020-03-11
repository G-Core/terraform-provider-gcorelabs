package fixture

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"

	th "gcloud/gcorecloud-go/testhelper"
	"gcloud/gcorecloud-go/testhelper/client"
)

func SetupHandler(t *testing.T, url, method, requestBody, responseBody string, status int) {
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, method)
		th.TestHeader(t, r, "X-Auth-AccessToken", client.TokenID)

		if requestBody != "" {
			th.TestJSONRequest(t, r, requestBody)
		}

		if responseBody != "" {
			w.Header().Add("Content-Type", "application/json")
		}

		w.WriteHeader(status)

		if responseBody != "" {
			_, err := fmt.Fprint(w, responseBody)
			if err != nil {
				log.Error(err)
			}
		}
	})
}

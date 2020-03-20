package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

var stackID = "stack"
var resourceName = "resource"

func prepareResourceTestURLParams(projectID int, regionID int, stackID, resourceName, action string) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks/%s/resources/%s/%s", projectID, regionID, stackID, resourceName, action)
}

func prepareMetadataTestURL(stackID, resourceName string) string {
	return prepareResourceTestURLParams(fake.ProjectID, fake.RegionID, stackID, resourceName, "metadata")
}

func prepareSignalTestURL(stackID, resourceName string) string {
	return prepareResourceTestURLParams(fake.ProjectID, fake.RegionID, stackID, resourceName, "signal")
}

func TestResourceMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(stackID, resourceName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")
	metadata, err := resources.Metadata(client, stackID, resourceName).Extract()
	require.NoError(t, err)
	require.Equal(t, Metadata, metadata)
}

func TestResourceSignal(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareSignalTestURL(stackID, resourceName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestBody(t, r, SignalRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	client := fake.ServiceTokenClient("heat", "v1")
	err := resources.Signal(client, stackID, resourceName, []byte(SignalRequest)).ExtractErr()
	require.NoError(t, err)
}

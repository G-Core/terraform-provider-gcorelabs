package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

var stackID = "stack"
var resourceName = "resource"

func prepareResourceActionTestURLParams(projectID int, regionID int, stackID, resourceName, action string) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks/%s/resources/%s/%s", projectID, regionID, stackID, resourceName, action)
}

func prepareResourceTestURLParams(projectID int, regionID int, stackID, resourceName string) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks/%s/resources/%s", projectID, regionID, stackID, resourceName)
}

func prepareResourcesTestURLParams(projectID int, regionID int, stackID string) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks/%s/resources", projectID, regionID, stackID)
}

func prepareMetadataTestURL(stackID, resourceName string) string {
	return prepareResourceActionTestURLParams(fake.ProjectID, fake.RegionID, stackID, resourceName, "metadata")
}

func prepareSignalTestURL(stackID, resourceName string) string {
	return prepareResourceActionTestURLParams(fake.ProjectID, fake.RegionID, stackID, resourceName, "signal")
}

func prepareGetResourceTestURL(stackID, resourceName string) string {
	return prepareResourceTestURLParams(fake.ProjectID, fake.RegionID, stackID, resourceName)
}

func prepareListResourcesTestURL(stackID string) string {
	return prepareResourcesTestURLParams(fake.ProjectID, fake.RegionID, stackID)
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

func TestGetResource(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetResourceTestURL(stackID, resourceName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}

	})

	client := fake.ServiceTokenClient("heat", "v1")
	resource, err := resources.Get(client, stackID, resourceName).Extract()
	require.NoError(t, err)
	require.Equal(t, StackResource1, *resource)

}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	url := prepareListResourcesTestURL(stackID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")
	count := 0

	err := resources.List(client, stackID, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := resources.ExtractResources(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, StackResourceList1, ct)
		require.Equal(t, ExpectedStackResourceList1, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	url := prepareListResourcesTestURL(stackID)
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")

	actual, err := resources.ListAll(client, stackID, nil)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, StackResourceList1, ct)
	require.Equal(t, ExpectedStackResourceList1, actual)
	th.AssertNoErr(t, err)

}

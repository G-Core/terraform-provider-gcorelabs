package testing

import (
	"encoding/json"
	"fmt"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/flavor/v1/flavors"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/flavors/%d/%d", projectID, regionID)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("flavors", "v1")
	count := 0

	opts := flavors.ListOpts{}

	err := flavors.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := flavors.ExtractFlavors(page)
		require.NoError(t, err)
		flavor := actual[0]
		require.Equal(t, Flavor1, flavor)
		require.Equal(t, ExpectedFlavorSlice, actual)

		list, err := json.Marshal(ExpectedFlavorSlice)
		require.NoError(t, err)

		require.JSONEq(t, ListSliceResponse, string(list))

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestFindByName(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("flavors", "v1")

	flavorID, err := flavors.IDFromName(client, Flavor1.FlavorName)
	require.NoError(t, err)
	require.Equal(t, Flavor1.FlavorID, flavorID)

	_, err = flavors.IDFromName(client, "X")
	require.Error(t, err)

}

func TestListMalformedResponse(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponseMalformedCurrency)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("flavors", "v1")

	_, err := flavors.List(client, nil).AllPages()
	require.Error(t, err)

}

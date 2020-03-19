package testing

import (
	"fmt"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/subnet/v1/subnets"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/subnets/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/subnets/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
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

	client := fake.ServiceTokenClient("subnets", "v1")
	count := 0

	opts := subnets.ListOpts{
		NetworkID: Subnet1.NetworkID,
	}

	err := subnets.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := subnets.ExtractSubnets(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Subnet1, ct)
		require.Equal(t, ExpectedSubnetSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")

	ct, err := subnets.Get(client, Subnet1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)

}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             true,
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
	}

	client := fake.ServiceTokenClient("subnets", "v1")
	tasks, err := subnets.Create(client, options).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Subnet1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")
	tasks, err := subnets.Delete(client, Subnet1.ID).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")

	opts := subnets.UpdateOpts{
		Name: Subnet1.Name,
	}

	ct, err := subnets.Update(client, Subnet1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
	require.Equal(t, Subnet1.Name, ct.Name)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)

}

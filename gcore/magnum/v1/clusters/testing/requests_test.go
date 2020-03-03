package testing

import (
	"fmt"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	fake "gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"gcloud/gcorecloud-go/pagination"
	th "gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/clusters", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/clusters/%s", projectID, regionID, id)
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
		_, err := fmt.Fprintf(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("magnum", "v1")
	count := 0

	err := clusters.List(client, clusters.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, ClusterList1, ct)
		require.Equal(t, ExpectedClusterSlice, actual)
		require.Nil(t, ct.HealthStatus)
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

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintf(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("magnum", "v1")

	ct, err := clusters.Get(client, Cluster1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, Cluster1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)
	require.Nil(t, ct.HealthStatus)

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

		_, err := fmt.Fprintf(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	timeout := 360

	options := clusters.CreateOpts{
		Name:              Cluster1.Name,
		ClusterTemplateId: Cluster1.ClusterTemplateID,
		NodeCount:         1,
		MasterCount:       1,
		KeyPair:           "keypair",
		FlavorId:          Cluster1.FlavorID,
		DiscoveryUrl:      nil,
		CreateTimeout:     &timeout,
		MasterFlavorId:    Cluster1.MasterFlavorID,
		Labels:            &map[string]string{},
		FixedSubnet:       nil,
	}

	client := fake.ServiceTokenClient("magnum", "v1")
	tasks, err := clusters.Create(client, options).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("magnum", "v1")
	res := clusters.Delete(client, Cluster1.UUID)
	th.AssertNoErr(t, res.Err)
}

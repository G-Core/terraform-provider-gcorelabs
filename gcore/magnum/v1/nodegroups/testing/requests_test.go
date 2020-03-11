package testing

import (
	"fmt"
	"gcloud/gcorecloud-go/gcore/magnum/v1/nodegroups"
	fake "gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"gcloud/gcorecloud-go/pagination"
	th "gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int, clusterID string) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/nodegroups/%s", projectID, regionID, clusterID)
}

func prepareGetTestURLParams(projectID int, regionID int, clusterID, id string) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/nodegroups/%s/%s", projectID, regionID, clusterID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID, NodeGroup1.ClusterID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, NodeGroup1.ClusterID, id)
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

	client := fake.ServiceTokenClient("magnum", "v1")
	count := 0

	err := nodegroups.List(client, NodeGroup1.ClusterID, nodegroups.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := nodegroups.ExtractClusterNodeGroups(page)
		require.NoError(t, err)
		ng1 := actual[0]
		ng2 := actual[1]
		require.Equal(t, NodeGroupList1, ng1)
		require.Equal(t, NodeGroupList2, ng2)
		require.Equal(t, ExpectedClusterNodeGroupSlice, actual)
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

	th.Mux.HandleFunc(prepareGetTestURL(NodeGroup1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("magnum", "v1")

	ct, err := nodegroups.Get(client, NodeGroup1.ClusterID, NodeGroup1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, NodeGroup1, *ct)
	th.CheckDeepEquals(t, &NodeGroup1, ct)

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

	size := 5

	options := nodegroups.CreateOpts{
		Name:             NodeGroup1.Name,
		FlavorID:         NodeGroup1.FlavorID,
		ImageID:          NodeGroup1.ImageID,
		NodeCount:        1,
		DockerVolumeSize: &size,
	}

	client := fake.ServiceTokenClient("magnum", "v1")
	tasks, err := nodegroups.Create(client, NodeGroup1.ClusterID, options).ExtractTasks()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	listenURL := prepareGetTestURL(NodeGroup1.UUID)
	th.Mux.HandleFunc(listenURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("magnum", "v1")
	tasks, err := nodegroups.Delete(client, NodeGroup1.ClusterID, NodeGroup1.UUID).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(NodeGroup1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UpdateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("magnum", "v1")

	maxNodeCount := 20

	options := nodegroups.UpdateOpts{
		MaxNodeCount: &maxNodeCount,
	}

	ct, err := nodegroups.Update(client, NodeGroup1.ClusterID, NodeGroup1.UUID, options).Extract()

	require.NoError(t, err)
	require.Equal(t, UpdatedNodeGroup1, *ct)
	th.CheckDeepEquals(t, &UpdatedNodeGroup1, ct)

}

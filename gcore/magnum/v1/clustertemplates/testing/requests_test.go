package testing

import (
	"fmt"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clustertemplates"
	fake "gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"gcloud/gcorecloud-go/pagination"
	th "gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/clustertemplates", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/magnum/%d/%d/clustertemplates/%s", projectID, regionID, id)
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

	client := fake.ServiceTokenClient("magnum", "v1")
	count := 0

	err := clustertemplates.List(client, clustertemplates.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clustertemplates.ExtractClusterTemplates(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, ExpectedClusterTemplateSlice, actual)
		require.Equal(t, ClusterTemplate1, ct)
		require.Equal(t, createdTime, ct.CreatedAt)
		require.Nil(t, ct.UpdatedAt)
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

	th.Mux.HandleFunc(prepareGetTestURL(ClusterTemplate1.UUID), func(w http.ResponseWriter, r *http.Request) {
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

	ct, err := clustertemplates.Get(client, ClusterTemplate1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, ClusterTemplate1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	th.CheckDeepEquals(t, &ClusterTemplate1, ct)
	require.Nil(t, ct.UpdatedAt)

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

	options := clustertemplates.CreateOpts{
		Name:             ClusterTemplate1.Name,
		ImageID:          ClusterTemplate1.ImageID,
		KeyPairID:        ClusterTemplate1.KeyPairID,
		DockerVolumeSize: ClusterTemplate1.DockerVolumeSize,
	}

	client := fake.ServiceTokenClient("magnum", "v1")
	ct, err := clustertemplates.Create(client, options).Extract()

	require.NoError(t, err)
	require.Equal(t, ClusterTemplate1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	th.CheckDeepEquals(t, &ClusterTemplate1, ct)
	require.Nil(t, ct.UpdatedAt)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(ClusterTemplate1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("magnum", "v1")
	res := clustertemplates.Delete(client, ClusterTemplate1.UUID)
	th.AssertNoErr(t, res.Err)
}

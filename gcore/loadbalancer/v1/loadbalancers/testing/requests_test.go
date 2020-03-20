package testing

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/loadbalancers"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, id)
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

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	count := 0

	err := loadbalancers.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := loadbalancers.ExtractLoadBalancers(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, LoadBalancer1, ct)
		require.Equal(t, ExpectedLoadBalancerSlice, actual)
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

	testURL := prepareGetTestURL(LoadBalancer1.ID)

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

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	ct, err := loadbalancers.Get(client, LoadBalancer1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, LoadBalancer1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)

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

	urlPath := "/"
	maxRetriesDown := 3
	memberWeight1 := 2
	memberWeight2 := 4
	instanceID1 := "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
	instanceID2 := "169942e0-9b53-42df-95ef-1a8b6525c2bd"

	options := loadbalancers.CreateOpts{
		Name: LoadBalancer1.Name,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:             "listener_name",
			ProtocolPort:     80,
			Protocol:         types.ProtocolTypeHTTP,
			Certificate:      nil,
			CertificateChain: nil,
			PrivateKey:       nil,
			Pools: []loadbalancers.CreatePoolOpts{{
				Name:     "pool_name",
				Protocol: types.ProtocolTypeHTTP,
				Members: []loadbalancers.CreatePoolMemberOpts{{
					InstanceID:   &instanceID1,
					Address:      net.ParseIP("192.168.1.101"),
					ProtocolPort: 8000,
					Weight:       &memberWeight1,
					SubnetID:     nil,
				}, {
					Address:      net.ParseIP("192.168.1.102"),
					ProtocolPort: 8000,
					Weight:       &memberWeight2,
					SubnetID:     nil,
					InstanceID:   &instanceID2,
				},
				},
				HealthMonitor: &loadbalancers.CreateHealthMonitorOpts{
					Type:           types.HealthMonitorTypeHTTP,
					Delay:          10,
					MaxRetries:     3,
					Timeout:        5,
					MaxRetriesDown: &maxRetriesDown,
					HTTPMethod:     types.HTTPMethodGET,
					URLPath:        &urlPath,
				},
				LoadBalancerAlgorithm: types.LoadBalancerAlgorithmRoundRobin,
				SessionPersistence:    nil,
			}},
		}},
		VipNetworkID: nil,
	}

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	tasks, err := loadbalancers.Create(client, options).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(LoadBalancer1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	tasks, err := loadbalancers.Delete(client, LoadBalancer1.ID).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(LoadBalancer1.ID)

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

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	opts := loadbalancers.UpdateOpts{
		Name: LoadBalancer1.Name,
	}

	ct, err := loadbalancers.Update(client, LoadBalancer1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, LoadBalancer1, *ct)
	require.Equal(t, LoadBalancer1.Name, ct.Name)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)

}

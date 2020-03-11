package testing

import (
	"fmt"
	"gcloud/gcorecloud-go/gcore/volume/v1/volumes"
	fake "gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"gcloud/gcorecloud-go/pagination"
	th "gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string { // nolint
	return fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareAttachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "attach")
}

func prepareDetachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "detach")
}

func prepareRetypeTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "retype")
}

func prepareExtendTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "extend")
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

	client := fake.ServiceTokenClient("volumes", "v1")
	count := 0

	opts := volumes.ListOpts{}

	err := volumes.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := volumes.ExtractVolumes(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Volume1, ct)
		require.Equal(t, ExpectedVolumeSlice, actual)
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

	testURL := prepareGetTestURL(Volume1.ID)

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

	client := fake.ServiceTokenClient("volumes", "v1")

	ct, err := volumes.Get(client, Volume1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Volume1, *ct)
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

	size := 10
	typeName := volumes.SsdHiIops
	instanceIDToAttachTo := "88f3e0bd-ca86-4cf7-be8b-dd2988e23c2d"

	options := volumes.CreateOpts{
		Source:               "new-volume",
		Name:                 "TestVM5 Ubuntu volume",
		Size:                 &size,
		TypeName:             &typeName,
		ImageID:              nil,
		SnapshotID:           nil,
		InstanceIDToAttachTo: &instanceIDToAttachTo,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	tasks, err := volumes.Create(client, options).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("volumes", "v1")

	opts := volumes.DeleteOpts{Snapshots: []string{"x", "y"}}

	tasks, err := volumes.Delete(client, Volume1.ID, opts).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestAttach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAttachTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AttachDetachRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Attach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
}

func TestDetach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDetachTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AttachDetachRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Detach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
}

func TestRetype(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRetypeTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RetypeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.VolumeTypePropertyOperationOpts{
		VolumeType: volumes.SsdHiIops,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Retype(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
}

func TestExtend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareExtendTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ExtendRequest)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ExtendResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("volumes", "v1")
	opts := volumes.SizePropertyOperationOpts{Size: 16}

	tasks, err := volumes.Extend(client, Volume1.ID, opts).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

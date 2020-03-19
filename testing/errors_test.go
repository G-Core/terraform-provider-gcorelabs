package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func TestGetResponseCode(t *testing.T) {
	respErr := gcorecloud.ErrUnexpectedResponseCode{
		URL:      "http://example.com",
		Method:   "GET",
		Expected: []int{200},
		Actual:   404,
		Body:     nil,
	}

	var err404 error = gcorecloud.ErrDefault404{ErrUnexpectedResponseCode: respErr}

	err, ok := err404.(gcorecloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 404)
}

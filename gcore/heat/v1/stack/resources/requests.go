package resources

import (
	"bytes"
	"gcloud/gcorecloud-go"
)

// Metadata retrieves metadata for heat resource
func Metadata(c *gcorecloud.ServiceClient, id, resource string) (r MetadataResult) {
	url := MetadataURL(c, id, resource)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Signal set heat resource status
func Signal(c *gcorecloud.ServiceClient, id, resource string, body []byte) (r SignalResult) {
	url := SignalURL(c, id, resource)
	_, r.Err = c.Post(url, nil, nil, &gcorecloud.RequestOpts{
		RawBody: bytes.NewReader(body),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	})
	return
}

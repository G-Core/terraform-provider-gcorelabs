package tokens

import (
	"gcloud/gcorecloud-go"
)

func processToken(c *gcorecloud.ServiceClient, opts gcorecloud.AuthOptionsBuilder, url string) (r TokenResult) {
	b := opts.ToMap()
	resp, err := c.Post(url, b, &r.Body, &gcorecloud.RequestOpts{})
	r.Err = err
	if resp != nil {
		r.Header = resp.Header
	}
	return
}

// Create authenticates and either generates a new token
func Create(c *gcorecloud.ServiceClient, opts gcorecloud.AuthOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, tokenURL(c))
}

// Refresh token with GCore API
func Refresh(c *gcorecloud.ServiceClient, opts gcorecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshURL(c))
}

// Refresh token with gcloud API
func RefreshGCloud(c *gcorecloud.ServiceClient, opts gcorecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshGCloudURL(c))
}

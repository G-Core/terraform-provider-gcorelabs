package stacks

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// List all stacks
func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return StackPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific heat stack.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

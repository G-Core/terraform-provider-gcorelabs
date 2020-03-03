package tasks

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/pagination"
)

// List returns a Pager which allows you to iterate over a collection of
// cluster templates. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return TaskPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster template based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

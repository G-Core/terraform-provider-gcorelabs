package keypairs

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return KeyPairPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific keypair based on its name or ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToKeyPairCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a keypair.
type CreateOpts struct {
	Name      string `json:"sshkey_name"`
	PublicKey string `json:"public_key"`
}

// ToKeyPairCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToKeyPairCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new keypair using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeyPairCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the keypair associated with it.
func Delete(c *gcorecloud.ServiceClient, keypairID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, keypairID), nil)
	return
}

// IDFromName is a convenience function that returns a keypair ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	pages, err := List(client).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractKeyPairs(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "keypairs"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "keypairs"}
	}
}

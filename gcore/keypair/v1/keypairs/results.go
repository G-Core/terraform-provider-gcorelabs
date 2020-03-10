package keypairs

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a keypair resource.
func (r commonResult) Extract() (*KeyPair, error) {
	var s KeyPair
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a KeyPair.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a KeyPair.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a KeyPair.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// KeyPair represents a keypair structure.
type KeyPair struct {
	Name        string  `json:"sshkey_name"`
	ID          string  `json:"sshkey_id"`
	Fingerprint string  `json:"fingerprint"`
	PublicKey   string  `json:"public_key"`
	PrivateKey  *string `json:"private_key"`
}

// KeyPairPage is the page returned by a pager when traversing over a
// collection of keypairs.
type KeyPairPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of keypairs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r KeyPairPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a KeyPairPage struct is empty.
func (r KeyPairPage) IsEmpty() (bool, error) {
	is, err := ExtractKeyPairs(r)
	return len(is) == 0, err
}

// ExtractKeyPair accepts a Page struct, specifically a KeyPairPage struct,
// and extracts the elements into a slice of KeyPair structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractKeyPairs(r pagination.Page) ([]KeyPair, error) {
	var s []KeyPair
	err := ExtractKeyPairsInto(r, &s)
	return s, err
}

func ExtractKeyPairsInto(r pagination.Page, v interface{}) error {
	return r.(KeyPairPage).Result.ExtractIntoSlicePtr(v, "results")
}

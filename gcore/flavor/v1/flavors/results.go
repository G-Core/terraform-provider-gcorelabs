package flavors

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"github.com/shopspring/decimal"
)

func init() { // nolint
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true
}

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a flavor resource.
func (r commonResult) Extract() (*Flavor, error) {
	var s Flavor
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Flavor.
type GetResult struct {
	commonResult
}

// Flavor represents a flavor structure.
type Flavor struct {
	FlavorID      string               `json:"flavor_id"`
	FlavorName    string               `json:"flavor_name"`
	PriceStatus   *string              `json:"price_status,omitempty"`
	CurrencyCode  *gcorecloud.Currency `json:"currency_code,omitempty"`
	PricePerHour  *decimal.Decimal     `json:"price_per_hour,omitempty"`
	PricePerMonth *decimal.Decimal     `json:"price_per_month,omitempty"`
	RAM           int                  `json:"ram"`
	VCPUS         int                  `json:"vcpus"`
}

// FlavorPage is the page returned by a pager when traversing over a
// collection of flavors.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of flavors has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FlavorPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FlavorPage struct is empty.
func (r FlavorPage) IsEmpty() (bool, error) {
	is, err := ExtractFlavors(r)
	return len(is) == 0, err
}

// ExtractFlavor accepts a Page struct, specifically a FlavorPage struct,
// and extracts the elements into a slice of Flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s []Flavor
	err := ExtractFlavorsInto(r, &s)
	return s, err
}

func ExtractFlavorsInto(r pagination.Page, v interface{}) error {
	return r.(FlavorPage).Result.ExtractIntoSlicePtr(v, "results")
}

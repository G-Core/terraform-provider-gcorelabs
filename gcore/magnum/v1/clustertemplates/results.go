package clustertemplates

import (
	"time"

	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster template resource.
func (r commonResult) Extract() (*ClusterTemplate, error) {
	var s ClusterTemplate
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Network.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Network.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Network.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gcorecloud.ErrResult
}

// ClusterTemplate represents a cluster template.
type ClusterTemplate struct {
	Labels              map[string]string `json:"labels"`
	FixedSubnet         string            `json:"fixed_subnet"`
	MasterFlavorID      string            `json:"master_flavor_id"`
	FlavorID            string            `json:"flavor_id"`
	NoProxy             string            `json:"no_proxy"`
	HttpsProxy          string            `json:"https_proxy"`
	HttpProxy           string            `json:"http_proxy"`
	TlsDisabled         bool              `json:"tls_disabled"`
	KeyPairID           string            `json:"keypair_id"`
	Public              bool              `json:"public"`
	DockerVolumeSize    int               `json:"docker_volume_size"`
	ServerType          string            `json:"server_type"`
	ExternalNetworkId   string            `json:"external_network_id"`
	ImageId             string            `json:"image_id"`
	VolumeDriver        string            `json:"volume_driver"`
	RegistryEnabled     bool              `json:"registry_enabled"`
	DockerStorageDriver string            `json:"docker_storage_driver"`
	Name                string            `json:"name"`
	NetworkDriver       string            `json:"network_driver"`
	FixedNetwork        string            `json:"fixed_network"`
	MasterLbEnabled     bool              `json:"master_lb_enabled"`
	DnsNameServer       string            `json:"dns_nameserver"`
	FloatingIpEnabled   bool              `json:"floating_ip_enabled"`
	Hidden              bool              `json:"hidden"`
	UUID                string            `json:"uuid"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           *time.Time        `json:"updated_at"`
	InsecureRegistry    string            `json:"insecure_registry"`
	Links               []gcorecloud.Link `json:"links"`
}

/*
func (r *ClusterTemplate) UnmarshalJSON(b []byte) error {
	var ct ClusterTemplate
	err := json.Unmarshal(b, &ct)
	if err != nil {
		return err
	}
	return nil
}
*/

// ClusterTemplatePage is the page returned by a pager when traversing over a
// collection of networks.
type ClusterTemplatePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of cluster templates has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterTemplatePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ClusterTemplate struct is empty.
func (r ClusterTemplatePage) IsEmpty() (bool, error) {
	is, err := ExtractClusterTemplates(r)
	return len(is) == 0, err
}

// ExtractClusterTemplates accepts a Page struct, specifically a ClusterTemplatePage struct,
// and extracts the elements into a slice of ClusterTemplate structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusterTemplates(r pagination.Page) ([]ClusterTemplate, error) {
	var s []ClusterTemplate
	err := ExtractClusterTemplatesInto(r, &s)
	return s, err
}

func ExtractClusterTemplatesInto(r pagination.Page, v interface{}) error {
	return r.(ClusterTemplatePage).Result.ExtractIntoSlicePtr(v, "clustertemplates")
}

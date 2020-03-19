package volumes

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the Create request
type CreateOptsBuilder interface {
	ToVolumeCreateMap() (map[string]interface{}, error)
}

// InstanceOperationOptsBuilder prepare data to proceed with Attach and Detach requests
type InstanceOperationOptsBuilder interface {
	ToVolumeInstanceOperationMap() (map[string]interface{}, error)
}

// PropertiesOperationOptsBuilder prepare data to proceed with Retype and Extend requests
type PropertiesOperationOptsBuilder interface {
	ToVolumePropertiesOperationMap() (map[string]interface{}, error)
}

// DeleteOptsBuilder allows extensions to add additional parameters to the Delete request
type DeleteOptsBuilder interface {
	ToVolumeDeleteQuery() (string, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToVolumeListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the cluster templates attributes you want to see returned. SortKey allows you to sort
// by a particular cluster templates attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	InstanceID *string `q:"instance_id"`
}

// CreateOpts represents options used to create a volume.
type CreateOpts struct {
	Source               VolumeSource `json:"source" required:"true"`
	Name                 string       `json:"name" required:"true"`
	Size                 *int         `json:"size,omitempty"`
	TypeName             *VolumeType  `json:"type_name,omitempty"`
	ImageID              *string      `json:"image_id,omitempty"`
	SnapshotID           *string      `json:"snapshot_id,omitempty"`
	InstanceIDToAttachTo *string      `json:"instance_id_to_attach_to,omitempty"`
}

// DeleteOpts allows set additional parameters during volume deletion.
type DeleteOpts struct {
	Snapshots []string `q:"snapshots" delimiter:"comma"`
}

// InstanceOperationOpts allows prepare data for Attach and Detach requests
type InstanceOperationOpts struct {
	InstanceID string `json:"instance_id,omitempty"`
}

type VolumeTypePropertyOperationOpts struct {
	VolumeType VolumeType `json:"volume_type,omitempty"`
}

type SizePropertyOperationOpts struct {
	Size int `json:"size,omitempty"`
}

// ToVolumeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ToVolumeDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToVolumeDeleteQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ToVolumeCreateMap builds a request body.
func (opts CreateOpts) ToVolumeCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToVolumeInstanceOperationMap builds a request body.
func (opts InstanceOperationOpts) ToVolumeInstanceOperationMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToVolumePropertiesOperationMap builds a request body.
func (opts VolumeTypePropertyOperationOpts) ToVolumePropertiesOperationMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToVolumePropertiesOperationMap builds a request body.
func (opts SizePropertyOperationOpts) ToVolumePropertiesOperationMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// List retrieves list of volumes
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VolumePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific volume based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Create accepts a CreateOpts struct and creates a new volume using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r TasksResult) {
	b, err := opts.ToVolumeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the volume associated with it.
func Delete(c *gcorecloud.ServiceClient, volumeID string, opts DeleteOptsBuilder) (r TasksResult) {
	url := deleteURL(c, volumeID)
	if opts != nil {
		query, err := opts.ToVolumeDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.DeleteWithResponse(url, &r.Body, nil)
	return
}

// Attach accepts a InstanceOperationOpts struct and attach volume to an instance.
func Attach(c *gcorecloud.ServiceClient, volumeID string, opts InstanceOperationOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumeInstanceOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(attachURL(c, volumeID), b, &r.Body, nil)
	return
}

// Detach accepts a InstanceOperationOpts struct and detach volume to an instance.
func Detach(c *gcorecloud.ServiceClient, volumeID string, opts InstanceOperationOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumeInstanceOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(detachURL(c, volumeID), b, &r.Body, nil)
	return
}

// Retype accepts a VolumeTypePropertyOperationOpts struct and retype volume.
func Retype(c *gcorecloud.ServiceClient, volumeID string, opts PropertiesOperationOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumePropertiesOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(retypeURL(c, volumeID), b, &r.Body, nil)
	return
}

// Extend accepts a VolumeTypePropertyOperationOpts struct and extend volume.
func Extend(c *gcorecloud.ServiceClient, volumeID string, opts PropertiesOperationOptsBuilder) (r TasksResult) {
	b, err := opts.ToVolumePropertiesOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(extendURL(c, volumeID), b, &r.Body, nil)
	return
}

// IDFromName is a convenience function that returns a volume ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	opts := ListOpts{}

	pages, err := List(client, opts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractVolumes(pages)
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
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "volumes"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "volumes"}
	}
}

package testing

import (
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources"
)

const MetadataResponse = `
{
	"some_key": "some_value",
	"some_other_key": "some_other_value"
}
`

const SignalRequest = `
{
	"some_key": "some_value",
	"some_other_key": "some_other_value"
}
`

var ListResponse = `
{
  "count": 20,
  "results": [
    {
      "resource_status": "UPDATE_COMPLETE",
      "resource_type": "file:///opt/stack/magnum/magnum/drivers/common/templates/lb_etcd.yaml",
      "updated_time": "2020-03-21T20:41:06+00:00",
      "physical_resource_id": "da7a7e75-28ff-4813-88a8-9ab583ac227f",
      "required_by": [
        "etcd_address_lb_switch",
        "kube_masters"
      ],
      "creation_time": "2020-03-17T21:26:07+00:00",
      "logical_resource_id": "etcd_lb",
      "resource_name": "etcd_lb",
      "resource_status_reason": "state changed"
    }
  ]	
}
`

var GetResponse = `
{
  "description": "",
  "resource_status": "UPDATE_COMPLETE",
  "attributes": {
    "True": null,
    "pool_id": "1a652ee9-1d09-4f33-a4b7-e81bf13da180",
    "address": "10.0.0.17"
  },
  "resource_type": "file:///opt/stack/magnum/magnum/drivers/common/templates/lb_etcd.yaml",
  "updated_time": "2020-03-21T20:41:06+00:00",
  "physical_resource_id": "da7a7e75-28ff-4813-88a8-9ab583ac227f",
  "required_by": [
    "etcd_address_lb_switch",
    "kube_masters"
  ],
  "creation_time": "2020-03-17T21:26:07+00:00",
  "logical_resource_id": "etcd_lb",
  "resource_name": "etcd_lb",
  "resource_status_reason": "state changed"
}
`

var (
	Metadata = map[string]interface{}{
		"some_key":       "some_value",
		"some_other_key": "some_other_value",
	}
	resourceStatusReason = "state changed"
	creationTime, _      = time.Parse(time.RFC3339, "2020-03-17T21:26:07+00:00")
	updatedTime, _       = time.Parse(time.RFC3339, "2020-03-21T20:41:06+00:00")
	StackResourceList1   = resources.ResourceList{
		CreationTime:         creationTime,
		UpdatedTime:          &updatedTime,
		LogicalResourceID:    "etcd_lb",
		PhysicalResourceID:   "da7a7e75-28ff-4813-88a8-9ab583ac227f",
		RequiredBy:           []string{"etcd_address_lb_switch", "kube_masters"},
		ResourceName:         "etcd_lb",
		ResourceStatus:       "UPDATE_COMPLETE",
		ResourceStatusReason: &resourceStatusReason,
		ResourceType:         "file:///opt/stack/magnum/magnum/drivers/common/templates/lb_etcd.yaml",
	}
	StackResource1 = resources.Resource{
		ResourceList: &StackResourceList1,
		Description:  "",
		Attributes: map[string]interface{}{
			"True":    nil,
			"pool_id": "1a652ee9-1d09-4f33-a4b7-e81bf13da180",
			"address": "10.0.0.17",
		},
	}
	ExpectedStackResourceList1 = []resources.ResourceList{StackResourceList1}
)

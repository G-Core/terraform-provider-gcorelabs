package testing

import (
	"net"
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/flavor/v1/flavors"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"
)

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
      "task_state": null,
      "instance_description": "Testing",
      "instance_name": "Testing",
      "status": "ACTIVE",
      "instance_created": "2019-07-11T06:58:48Z",
      "vm_state": "active",
      "volumes": [
        {
          "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
          "delete_on_termination": false
        }
      ],
      "security_groups": [
        {
          "name": "default"
        }
      ],
      "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
      "task_id": "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02",
      "creator_task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
      "addresses": {
        "net1": [
          {
            "type": "fixed",
            "addr": "10.0.0.17"
          },
          {
            "type": "floating",
            "addr": "92.38.157.215"
          }
        ],
        "net2": [
          {
            "type": "fixed",
            "addr": "192.168.68.68"
          }
        ]
      },
      "metadata": {
        "os_distro": "centos",
        "os_version": "1711-x64",
        "image_name": "cirros-0.3.5-x86_64-disk",
        "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
        "snapshot_name": "test_snapshot",
        "snapshot_id": "c286cd13-fba9-4302-9cdb-4351a05a56ea",
        "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
      },
      "flavor": {
        "flavor_name": "g1s-shared-1-0.5",
        "disk": 0,
        "flavor_id": "g1s-shared-1-0.5",
        "vcpus": 1,
        "ram": 512
      },
      "project_id": 1,
      "region_id": 1,
	  "region": "RegionOne"	
    }
  ]
}
`

const GetResponse = `
{
  "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
  "task_state": null,
  "instance_description": "Testing",
  "instance_name": "Testing",
  "status": "ACTIVE",
  "instance_created": "2019-07-11T06:58:48Z",
  "vm_state": "active",
  "volumes": [
    {
      "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
      "delete_on_termination": false
    }
  ],
  "security_groups": [
    {
      "name": "default"
    }
  ],
  "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
  "task_id": "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02",
  "creator_task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
  "addresses": {
    "net1": [
      {
        "type": "fixed",
        "addr": "10.0.0.17"
      },
      {
        "type": "floating",
        "addr": "92.38.157.215"
      }
    ],
    "net2": [
      {
        "type": "fixed",
        "addr": "192.168.68.68"
      }
    ]
  },
  "metadata": {
    "os_distro": "centos",
    "os_version": "1711-x64",
    "image_name": "cirros-0.3.5-x86_64-disk",
    "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
    "snapshot_name": "test_snapshot",
    "snapshot_id": "c286cd13-fba9-4302-9cdb-4351a05a56ea",
    "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
  },
  "flavor": {
    "flavor_name": "g1s-shared-1-0.5",
    "disk": 0,
    "flavor_id": "g1s-shared-1-0.5",
    "vcpus": 1,
    "ram": 512
  },
  "project_id": 1,
  "region_id": 1,
  "region": "RegionOne"	
}
`

var (
	ip1                 = net.ParseIP("10.0.0.17")
	ip2                 = net.ParseIP("92.38.157.215")
	ip3                 = net.ParseIP("192.168.68.68")
	tm, _               = time.Parse(gcorecloud.RFC3339ZZ, "2019-07-11T06:58:48Z")
	createdTime         = gcorecloud.JSONRFC3339ZZ{Time: tm}
	instanceID          = "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
	instanceName        = "Testing"
	instanceDescription = "Testing"
	taskID              = "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02"
	creatorTaskID       = "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"

	Instance1 = instances.Instance{
		ID:          instanceID,
		Name:        instanceName,
		Description: instanceDescription,
		CreatedAt:   createdTime,
		Status:      "ACTIVE",
		VMState:     "active",
		TaskState:   nil,
		Flavor: flavors.Flavor{
			FlavorID:   "g1s-shared-1-0.5",
			FlavorName: "g1s-shared-1-0.5",
			RAM:        512,
			VCPUS:      1,
		},
		Metadata: map[string]interface{}{
			"os_distro":     "centos",
			"os_version":    "1711-x64",
			"image_name":    "cirros-0.3.5-x86_64-disk",
			"image_id":      "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
			"snapshot_name": "test_snapshot",
			"snapshot_id":   "c286cd13-fba9-4302-9cdb-4351a05a56ea",
			"task_id":       "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
		},
		Volumes: []instances.InstanceVolume{{
			ID:                  "28bfe198-a003-4283-8dca-ab5da4a71b62",
			DeleteOnTermination: false,
		}},
		Addresses: map[string][]instances.InstanceAddress{
			"net1": {{
				Type:    types.AddressTypeFixed,
				Address: ip1,
			},
				{
					Type:    types.AddressTypeFloating,
					Address: ip2,
				},
			},
			"net2": {{
				Type:    types.AddressTypeFixed,
				Address: ip3,
			}},
		},
		SecurityGroups: []types.ItemName{{
			Name: "default",
		}},
		CreatorTaskID: &creatorTaskID,
		TaskID:        &taskID,
		ProjectID:     1,
		RegionID:      1,
		Region:        "RegionOne",
	}
	ExpectedInstancesSlice = []instances.Instance{Instance1}
)

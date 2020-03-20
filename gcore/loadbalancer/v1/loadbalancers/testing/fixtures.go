package testing

import (
	"net"
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/loadbalancers"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "region": "RegionOne",
      "created_at": "2020-01-24T13:57:12+0000",
      "name": "lb",
      "id": "e8ab1be4-1521-4266-be69-28dad4148a30",
      "provisioning_status": "ACTIVE",
      "updated_at": "2020-01-24T13:57:35+0000",
      "listeners": [
        {
          "id": "43658ea9-54bd-4807-90b1-925921c9a0d1"
        }
      ],
      "task_id": null,
      "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
      "vip_address": "5.5.5.5",
      "operating_status": "ONLINE",
      "project_id": 1,
      "region_id": 1
    }
  ]
}
`

const GetResponse = `
{
  "region": "RegionOne",
  "created_at": "2020-01-24T13:57:12+0000",
  "name": "lb",
  "id": "e8ab1be4-1521-4266-be69-28dad4148a30",
  "provisioning_status": "ACTIVE",
  "updated_at": "2020-01-24T13:57:35+0000",
  "listeners": [
    {
      "id": "43658ea9-54bd-4807-90b1-925921c9a0d1"
    }
  ],
  "task_id": null,
  "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
  "vip_address": "5.5.5.5",
  "operating_status": "ONLINE",
  "project_id": 1,
  "region_id": 1
}
`

const CreateRequest = `
{
  "name": "lb",
  "listeners": [
    {
      "name": "listener_name",
      "protocol": "HTTP",
      "protocol_port": 80,
      "pools": [
        {
          "name": "pool_name",
          "protocol": "HTTP",
          "members": [
            {
              "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
              "address": "192.168.1.101",
              "weight": 2,
              "protocol_port": 8000
            },
            {
              "instance_id": "169942e0-9b53-42df-95ef-1a8b6525c2bd",
              "address": "192.168.1.102",
              "weight": 4,
              "protocol_port": 8000
            }
          ],
          "healthmonitor": {
            "type": "HTTP",
            "delay": 10,
            "max_retries": 3,
            "timeout": 5,
            "max_retries_down": 3,
            "http_method": "GET",
            "url_path": "/"
          },
          "lb_algorithm": "ROUND_ROBIN"
        }
      ]
    }
  ]
}
`

const UpdateRequest = `
{
	"name": "lb"
}	
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	createdTimeString    = "2020-01-24T13:57:12+0000"
	updatedTimeString    = "2020-01-24T13:57:35+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
	updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
	updatedTime          = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
	creatorTaskID        = "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22"

	LoadBalancer1 = loadbalancers.LoadBalancer{
		Name:               "lb",
		ID:                 "e8ab1be4-1521-4266-be69-28dad4148a30",
		ProvisioningStatus: types.ProvisioningStatusActive,
		OperationStatus:    types.OperatingStatusOnline,
		VipAddress:         net.ParseIP("5.5.5.5"),
		Listeners: []types.ItemID{{
			ID: "43658ea9-54bd-4807-90b1-925921c9a0d1",
		}},
		CreatorTaskID: &creatorTaskID,
		TaskID:        nil,
		CreatedAt:     createdTime,
		UpdatedAt:     &updatedTime,
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
		Region:        "RegionOne",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedLoadBalancerSlice = []loadbalancers.LoadBalancer{LoadBalancer1}
)

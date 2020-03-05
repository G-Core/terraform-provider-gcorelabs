package testing

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/network/v1/networks"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	fake "gcloud/gcorecloud-go/testhelper/client"
	"time"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "creator_task_id": null,
      "region": "RegionOne",
      "name": "private",
      "mtu": 1450,
      "id": "e7944e55-f957-413d-aa56-fdc876543113",
      "updated_at": "2020-03-05T12:03:25+0000",
      "created_at": "2020-03-05T12:03:24+0000",
      "task_id": null,
      "region_id": 1,
      "shared": false,
      "subnets": [
        "3730b4d3-9337-4a60-a35e-7e1620aabe6f"
      ],
      "external": false,
      "project_id": 1
	}
  ]
}
`

const GetResponse = `
{
  "creator_task_id": null,
  "region": "RegionOne",
  "name": "private",
  "mtu": 1450,
  "id": "e7944e55-f957-413d-aa56-fdc876543113",
  "updated_at": "2020-03-05T12:03:25+0000",
  "created_at": "2020-03-05T12:03:24+0000",
  "task_id": null,
  "region_id": 1,
  "shared": false,
  "subnets": [
    "3730b4d3-9337-4a60-a35e-7e1620aabe6f"
  ],
  "external": false,
  "project_id": 1
}
`

const CreateRequest = `
{
	"name": "private",
	"mtu": 1450,
	"create_router": true
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

var createdTimeString = "2020-03-05T12:03:24+0000"
var updatedTimeString = "2020-03-05T12:03:25+0000"
var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339Z(createdTimeParsed)
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z(updatedTimeParsed)

var (
	Network1 = networks.Network{
		Name: "private",
		ID:   "e7944e55-f957-413d-aa56-fdc876543113",
		Subnets: []string{
			"3730b4d3-9337-4a60-a35e-7e1620aabe6f",
		},
		MTU:       1450,
		CreatedAt: createdTime,
		UpdatedAt: &updatedTime,
		External:  false,
		Default:   false,
		Shared:    false,
		ProjectID: fake.ProjectID,
		RegionID:  fake.RegionID,
		Region:    "RegionOne",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedNetworkSlice = []networks.Network{Network1}
)

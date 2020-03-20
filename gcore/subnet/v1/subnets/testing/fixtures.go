package testing

import (
	"net"
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/subnet/v1/subnets"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "id": "e7944e55-f957-413d-aa56-fdc876543113",
      "name": "subnet",
      "ip_version": 4,
      "enable_dhcp": true,
      "cidr": "192.168.10.0/24",
      "created_at": "2020-03-05T12:03:24+0000",
      "updated_at": "2020-03-05T12:03:25+0000",
	  "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
	  "task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
	  "creator_task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
	  "region": "RegionOne",
      "project_id": 1,
      "region_id": 1
    }
  ]
}
`

const GetResponse = `
{
  "id": "e7944e55-f957-413d-aa56-fdc876543113",
  "name": "subnet",
  "ip_version": 4,
  "enable_dhcp": true,
  "cidr": "192.168.10.0/24",
  "created_at": "2020-03-05T12:03:24+0000",
  "updated_at": "2020-03-05T12:03:25+0000",
  "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
  "creator_task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
  "region": "RegionOne",
  "project_id": 1,
  "region_id": 1
}
`

const CreateRequest = `
{
  "name": "subnet",
  "enable_dhcp": true,
  "cidr": "192.168.10.0/24",
  "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "connect_to_network_router": true
}
`

const UpdateRequest = `
{
	"name": "subnet"
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
var createdTime = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
var _, nt, _ = net.ParseCIDR("192.168.10.0/24")
var cidr = gcorecloud.CIDR{IPNet: *nt}
var taskID = "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"

var (
	Subnet1 = subnets.Subnet{
		ID:            "e7944e55-f957-413d-aa56-fdc876543113",
		Name:          "subnet",
		IPVersion:     4,
		EnableDHCP:    true,
		CIDR:          cidr,
		CreatedAt:     createdTime,
		UpdatedAt:     updatedTime,
		NetworkID:     "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		TaskID:        &taskID,
		CreatorTaskID: &taskID,
		Region:        "RegionOne",
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedSubnetSlice = []subnets.Subnet{Subnet1}
)

package testing

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"time"
)

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "loadbalancers": [
        "{'id': '79943b39-5e67-47e1-8878-85044b39667a'}"
      ],
      "session_persistence": null,
      "name": "lbaas_test_pool",
      "id": "9fccf0a3-c0de-441d-9afd-2b9b58b08b9f",
      "provisioning_status": "ACTIVE",
      "protocol": "TCP",
      "members": [
        {
          "address": "192.168.13.9",
          "id": "65f4e0eb-7846-490e-b44d-726c8baf3c25",
          "weight": 1,
          "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
          "protocol_port": 80
        },
        {
          "address": "192.168.13.8",
          "id": "f6a9c5dd-f8cc-448d-8e57-81de69d127cb",
          "weight": 1,
          "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
          "protocol_port": 80
        }
      ],
      "lb_algorithm": "ROUND_ROBIN",
      "task_id": null,
      "creator_task_id": "d8334c12-2881-4c4a-84ad-1b21fea73ad1",
      "listeners": [
        "{'id': 'c63341da-ea44-4027-bbf6-1f1939c783da'}"
      ],
      "operating_status": "ONLINE"
    }
  ]
}
`

const GetResponse = `
{
  "loadbalancers": [
    "{'id': '79943b39-5e67-47e1-8878-85044b39667a'}"
  ],
  "session_persistence": null,
  "name": "lbaas_test_pool",
  "id": "9fccf0a3-c0de-441d-9afd-2b9b58b08b9f",
  "provisioning_status": "ACTIVE",
  "protocol": "TCP",
  "members": [
    {
      "address": "192.168.13.9",
      "id": "65f4e0eb-7846-490e-b44d-726c8baf3c25",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    },
    {
      "address": "192.168.13.8",
      "id": "f6a9c5dd-f8cc-448d-8e57-81de69d127cb",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    }
  ],
  "lb_algorithm": "ROUND_ROBIN",
  "task_id": null,
  "creator_task_id": "d8334c12-2881-4c4a-84ad-1b21fea73ad1",
  "listeners": [
    "{'id': 'c63341da-ea44-4027-bbf6-1f1939c783da'}"
  ],
  "operating_status": "ONLINE"
}
`

const CreateRequest = `
{
  "loadbalancers": [
    "{'id': '79943b39-5e67-47e1-8878-85044b39667a'}"
  ],
  "session_persistence": null,
  "name": "lbaas_test_pool",
  "id": "9fccf0a3-c0de-441d-9afd-2b9b58b08b9f",
  "provisioning_status": "ACTIVE",
  "protocol": "TCP",
  "members": [
    {
      "address": "192.168.13.9",
      "id": "65f4e0eb-7846-490e-b44d-726c8baf3c25",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    },
    {
      "address": "192.168.13.8",
      "id": "f6a9c5dd-f8cc-448d-8e57-81de69d127cb",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    }
  ],
  "lb_algorithm": "ROUND_ROBIN",
  "task_id": null,
  "creator_task_id": "d8334c12-2881-4c4a-84ad-1b21fea73ad1",
  "listeners": [
    "{'id': 'c63341da-ea44-4027-bbf6-1f1939c783da'}"
  ],
  "operating_status": "ONLINE"
}
`

const CreatePoolMemberRequest = `
{
  "id": "string",
  "address": "string",
  "protocol_port": 0,
  "weight": 0,
  "subnet_id": "string",
  "instance_id": "string"
}
`

const UpdateRequest = `
{
	"name": "private"
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

var (
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)

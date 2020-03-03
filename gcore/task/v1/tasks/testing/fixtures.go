package testing

import (
	"encoding/json"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"time"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "state": "NEW",
      "updated_on": null,
      "task_type": "create_vm",
      "error": null,
      "client_id": 2,
      "user_client_id": 2,
      "data": {
        "name": "cirroz1",
        "reservation_id": "01d4925e-f5db-414a-9808-74e08aa4a741",
        "block_device_config": [
          {
            "source": "image",
            "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
            "name": "volume for vm1",
            "type_name": "standard",
            "size": 1,
            "boot_index": 0
          }
        ],
        "network_config": null,
        "security_groups": null,
        "keypair_name": null,
        "flavor_name": "g1s-shared-1-0.5"
      },
      "request_id": null,
      "id": "26986bc0-748a-4448-b836-0a2aa465eb06",
      "user_id": 3,
      "finished_on": null,
      "project_id": 444,
      "created_on": "2019-06-25T08:42:42",
      "created_resources": null
    }
  ]
}
`

const GetResponse = `
{
  "state": "NEW",
  "updated_on": null,
  "task_type": "create_vm",
  "error": null,
  "client_id": 2,
  "user_client_id": 2,
  "data": {
    "name": "cirroz1",
    "reservation_id": "01d4925e-f5db-414a-9808-74e08aa4a741",
    "block_device_config": [
      {
        "source": "image",
        "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
        "name": "volume for vm1",
        "type_name": "standard",
        "size": 1,
        "boot_index": 0
      }
    ],
    "network_config": null,
    "security_groups": null,
    "keypair_name": null,
    "flavor_name": "g1s-shared-1-0.5"
  },
  "request_id": null,
  "id": "26986bc0-748a-4448-b836-0a2aa465eb06",
  "user_id": 3,
  "finished_on": null,
  "project_id": 444,
  "created_on": "2019-06-25T08:42:42",
  "created_resources": null
}
`

var (
	createdTimeString    = "2019-06-25T08:42:42"
	taskID               = "26986bc0-748a-4448-b836-0a2aa465eb06"
	taskType             = "create_vm"
	projectID            = 444
	clientID             = 2
	userClientID         = 2
	userID               = 3
	createdTimeAsTime, _ = time.Parse(gcorecloud.RFC3339NoZ, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339NoZ(createdTimeAsTime)
	taskData             = []byte(`
		{
			"name": "cirroz1",
			"reservation_id": "01d4925e-f5db-414a-9808-74e08aa4a741",
			"block_device_config": [
			  {
				"source": "image",
				"image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
				"name": "volume for vm1",
				"type_name": "standard",
				"size": 1,
				"boot_index": 0
			  }
			],
			"network_config": null,
			"security_groups": null,
			"keypair_name": null,
			"flavor_name": "g1s-shared-1-0.5"
		}			
	`)
	dataMap = map[string]interface{}{}
	_       = json.Unmarshal(taskData, &dataMap)

	Task1 = tasks.Task{
		ID:               taskID,
		TaskType:         taskType,
		ProjectID:        projectID,
		ClientID:         clientID,
		RegionID:         nil,
		UserID:           userID,
		UserClientID:     userClientID,
		State:            tasks.TaskStateNew,
		CreatedOn:        gcorecloud.JSONRFC3339NoZ(createdTime),
		UpdatedOn:        nil,
		FinishedOn:       nil,
		AcknowledgedAt:   nil,
		AcknowledgedBy:   nil,
		CreatedResources: nil,
		RequestID:        nil,
		Error:            nil,
		Data:             &dataMap,
	}

	ExpectedTasks = tasks.Tasks{Task1}
)

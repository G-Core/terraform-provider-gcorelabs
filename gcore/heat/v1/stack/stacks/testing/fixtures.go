package testing

import (
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
	{
	  "id": "94b6fff4-16ac-4049-a4e3-5323c1ab6060",
	  "updated_time": "2020-03-21T20:41:21+00:00",
	  "stack_user_project_id": "39c030e9c44242e7b164c32a415f90c7",
	  "creation_time": "2020-03-17T21:59:54+00:00",
	  "deletion_time": null,
	  "stack_owner": null,
	  "stack_status_reason": "Stack UPDATE completed successfully",
	  "tags": null,
	  "description": "No description",
	  "stack_name": "test-53yd6kruws7j-kube_minions-i75dlbg4rtuk",
	  "parent": "9bea2b3f-a687-4a64-9a25-e1c7ed592d08",
	  "stack_status": "UPDATE_COMPLETE"
	}
  ]
}
`

const GetResponse = `
{
  "id": "94b6fff4-16ac-4049-a4e3-5323c1ab6060",
  "updated_time": "2020-03-21T20:41:21+00:00",
  "stack_user_project_id": "39c030e9c44242e7b164c32a415f90c7",
  "disable_rollback": true,
  "creation_time": "2020-03-17T21:59:54+00:00",
  "deletion_time": null,
  "timeout_mins": 7200,
  "stack_owner": null,
  "stack_status_reason": "Stack UPDATE completed successfully",
  "tags": null,
  "outputs": [
    {
      "output_value": [
        null,
        null
      ],
      "output_key": "kube_minion_external_ip",
      "description": "No description given"
    },
    {
      "output_value": {
        "1": "e4497e3d-6dc1-4bb7-ab2c-95f2d9e0ddb6",
        "0": "96dc3da2-d305-403e-af04-84e6e3718794"
      },
      "output_key": "refs_map",
      "description": "No description given"
    },
    {
      "output_value": [
        "10.0.0.31",
        "10.0.0.13"
      ],
      "output_key": "kube_minion_ip",
      "description": "No description given"
    }
  ],
  "parameters": {
    "OS::project_id": "46beed3938e6474390b530fefd6173d2",
    "OS::stack_id": "94b6fff4-16ac-4049-a4e3-5323c1ab6060",
    "OS::stack_name": "test-53yd6kruws7j-kube_minions-i75dlbg4rtuk"
  },
  "notification_topics": [],
  "description": "No description",
  "stack_name": "test-53yd6kruws7j-kube_minions-i75dlbg4rtuk",
  "capabilities": [],
  "template_description": "No description",
  "parent": "9bea2b3f-a687-4a64-9a25-e1c7ed592d08",
  "stack_status": "UPDATE_COMPLETE"
}
`

var (
	creationTime, _     = time.Parse(time.RFC3339, "2020-03-17T21:59:54+00:00")
	updatedTime, _      = time.Parse(time.RFC3339, "2020-03-21T20:41:21+00:00")
	parent              = "9bea2b3f-a687-4a64-9a25-e1c7ed592d08"
	stackStatusReason   = "Stack UPDATE completed successfully"
	templateDescription = "No description"

	StackList1 = stacks.StackList{
		CreationTime:       creationTime,
		DeletionTime:       nil,
		UpdatedTime:        &updatedTime,
		Description:        "No description",
		ID:                 "94b6fff4-16ac-4049-a4e3-5323c1ab6060",
		Parent:             &parent,
		StackName:          "test-53yd6kruws7j-kube_minions-i75dlbg4rtuk",
		StackOwner:         nil,
		StackStatus:        "UPDATE_COMPLETE",
		StackStatusReason:  &stackStatusReason,
		StackUserProjectID: "39c030e9c44242e7b164c32a415f90c7",
		Tags:               nil,
	}

	Stack1 = stacks.Stack{
		StackList:           &StackList1,
		Capabilities:        []string{},
		DisableRollback:     true,
		NotificationTopics:  []string{},
		TemplateDescription: &templateDescription,
		TimeoutMinutes:      7200,
		Outputs: []map[string]interface{}{
			{
				"output_value": []interface{}{nil, nil},
				"output_key":   "kube_minion_external_ip",
				"description":  "No description given",
			},
			{
				"output_value": map[string]interface{}{
					"1": "e4497e3d-6dc1-4bb7-ab2c-95f2d9e0ddb6",
					"0": "96dc3da2-d305-403e-af04-84e6e3718794",
				},
				"output_key":  "refs_map",
				"description": "No description given",
			},
			{
				"output_value": []interface{}{
					"10.0.0.31",
					"10.0.0.13",
				},
				"output_key":  "kube_minion_ip",
				"description": "No description given",
			},
		},
		Parameters: map[string]interface{}{
			"OS::project_id": "46beed3938e6474390b530fefd6173d2",
			"OS::stack_id":   "94b6fff4-16ac-4049-a4e3-5323c1ab6060",
			"OS::stack_name": "test-53yd6kruws7j-kube_minions-i75dlbg4rtuk",
		},
	}

	ExpectedStackList1 = []stacks.StackList{StackList1}
)

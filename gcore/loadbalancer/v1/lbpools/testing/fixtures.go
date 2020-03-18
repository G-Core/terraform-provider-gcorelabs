package testing

import (
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/lbpools"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"net"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "loadbalancers": [
		{"id": "79943b39-5e67-47e1-8878-85044b39667a"}
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
        {"id": "c63341da-ea44-4027-bbf6-1f1939c783da"}
      ],
      "operating_status": "ONLINE"
    }
  ]
}
`

const GetResponse = `
{
  "loadbalancers": [
    {"id": "79943b39-5e67-47e1-8878-85044b39667a"}
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
	{"id": "c63341da-ea44-4027-bbf6-1f1939c783da"}
  ],
  "operating_status": "ONLINE"
}
`

const CreateRequest = `
{
  "loadbalancer_id": "79943b39-5e67-47e1-8878-85044b39667a",
  "name": "lbaas_test_pool",
  "protocol": "TCP",
  "members": [
    {
      "address": "192.168.13.9",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    },
    {
      "address": "192.168.13.8",
      "weight": 1,
      "subnet_id": "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d",
      "protocol_port": 80
    }
  ],
  "lb_algorithm": "ROUND_ROBIN",
  "listener_id": "c63341da-ea44-4027-bbf6-1f1939c783da"
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
	"name": "lbaas_test_pool"
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

const UpdateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	ip1            = net.ParseIP("192.168.13.9")
	ip2            = net.ParseIP("192.168.13.8")
	subnetID       = "c864873b-8d9b-4d29-8cce-bf0bdfdaa74d"
	LoadBalancerID = "79943b39-5e67-47e1-8878-85044b39667a"
	ListenerID     = "c63341da-ea44-4027-bbf6-1f1939c783da"
	creatorTaskID  = "d8334c12-2881-4c4a-84ad-1b21fea73ad1"
	width          = 1
	protocolPort   = 80
	Member1        = lbpools.PoolMember{
		Address:      &ip1,
		ID:           "65f4e0eb-7846-490e-b44d-726c8baf3c25",
		Weight:       &width,
		SubnetID:     &subnetID,
		InstanceID:   nil,
		ProtocolPort: &protocolPort,
	}
	Member2 = lbpools.PoolMember{
		Address:      &ip2,
		ID:           "f6a9c5dd-f8cc-448d-8e57-81de69d127cb",
		Weight:       &width,
		SubnetID:     &subnetID,
		InstanceID:   nil,
		ProtocolPort: &protocolPort,
	}
	LBPool1 = lbpools.Pool{
		LoadBalancers: []types.ItemID{
			{ID: LoadBalancerID},
		},
		Listeners: []types.ItemID{
			{ID: ListenerID},
		},
		SessionPersistence:    nil,
		LoadBalancerAlgorithm: types.LoadBalancerAlgorithmRoundRobin,
		Name:                  "lbaas_test_pool",
		ID:                    "9fccf0a3-c0de-441d-9afd-2b9b58b08b9f",
		Protocol:              types.ProtocolTypeTCP,
		Members: []lbpools.PoolMember{
			Member1,
			Member2,
		},
		ProvisioningStatus: types.ProvisioningStatusActive,
		OperatingStatus:    types.OperatingStatusOnline,
		CreatorTaskID:      &creatorTaskID,
		TaskID:             nil,
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
	ExpectedLBPoolsSlice = []lbpools.Pool{LBPool1}
)

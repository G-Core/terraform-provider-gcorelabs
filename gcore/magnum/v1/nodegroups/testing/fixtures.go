package testing

import (
	"gcloud/gcorecloud-go/gcore/magnum/v1/nodegroups"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"net"
)

const ListResponse = `
{
  "count": 2,	
  "results": [
    {
      "status": "CREATE_IN_PROGRESS",
      "name": "default-master",
      "node_count": 1,
      "uuid": "3eda6b46-58d9-4abc-8a11-6045b791a35b",
      "image_id": "fedora-coreos",
      "flavor_id": "g1-standard-1-2",
      "role": "master"
    },
    {
      "status": "CREATE_IN_PROGRESS",
      "name": "default-worker",
      "node_count": 1,
      "uuid": "467fc654-2d1c-48e8-9d59-489e8fcf8c17",
      "image_id": "fedora-coreos",
      "flavor_id": "g1-standard-1-2",
      "role": "worker"
    }
  ]
}
`

const GetResponse = `
{
  "project_id": "46beed3938e6474390b530fefd6173d2",
  "status": "CREATE_IN_PROGRESS",
  "cluster_id": "c94f38cc-dc78-4715-8939-68de082bd5e3",
  "docker_volume_size": 10,
  "min_node_count": 1,
  "name": "default-master",
  "stack_id": "ccff14f2-4585-407d-aea5-8581aa9a292f",
  "status_reason": null,
  "node_count": 1,
  "uuid": "3eda6b46-58d9-4abc-8a11-6045b791a35b",
  "image_id": "fedora-coreos",
  "node_addresses": [],
  "flavor_id": "g1-standard-1-2",
  "max_node_count": null,
  "labels": {
    "gcloud_project_id": "1",
    "gcloud_region_id": "1",
    "gcloud_access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4OTM3Nzc5LCJqdGkiOiI2NDk1MzFjMDNmYmU4NDczY2RiZGY2MjJkMWNmN2YzMCIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6MSwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxLCJ1c2VyX2lkIjoxLCJpc19hZG1pbiI6ZmFsc2V9.kdx_kpXF_Z_aCPDP5C-wIVLgv-mW9SSafJ_u6x7XO_k",
    "gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5OTExMzc3OSwianRpIjoiNjQ5NTMxYzAzZmJlODQ3M2NkYmRmNjIyZDFjZjdmMzAiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOjEsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MSwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.F_KTIGt1uvEaKnb8ZziI7Xca1o7Vcwj7752qbfb3Otg"
  },
  "role": "master",
  "is_default": true
}
`

const UpdateResponse = `
{
  "project_id": "46beed3938e6474390b530fefd6173d2",
  "status": "CREATE_IN_PROGRESS",
  "cluster_id": "c94f38cc-dc78-4715-8939-68de082bd5e3",
  "docker_volume_size": 10,
  "min_node_count": 1,
  "name": "default-master",
  "stack_id": "ccff14f2-4585-407d-aea5-8581aa9a292f",
  "status_reason": null,
  "node_count": 1,
  "uuid": "3eda6b46-58d9-4abc-8a11-6045b791a35b",
  "image_id": "fedora-coreos",
  "node_addresses": [],
  "flavor_id": "g1-standard-1-2",
  "max_node_count": 20,
  "labels": {
    "gcloud_project_id": "1",
    "gcloud_region_id": "1",
    "gcloud_access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4OTM3Nzc5LCJqdGkiOiI2NDk1MzFjMDNmYmU4NDczY2RiZGY2MjJkMWNmN2YzMCIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6MSwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxLCJ1c2VyX2lkIjoxLCJpc19hZG1pbiI6ZmFsc2V9.kdx_kpXF_Z_aCPDP5C-wIVLgv-mW9SSafJ_u6x7XO_k",
    "gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5OTExMzc3OSwianRpIjoiNjQ5NTMxYzAzZmJlODQ3M2NkYmRmNjIyZDFjZjdmMzAiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOjEsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MSwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.F_KTIGt1uvEaKnb8ZziI7Xca1o7Vcwj7752qbfb3Otg"
  },
  "role": "master",
  "is_default": true
}
`

const CreateRequest = `
{
  "docker_volume_size": 5,
  "name": "default-master",
  "node_count": 1,
  "flavor_id": "g1-standard-1-2",
  "image_id": "fedora-coreos"
}
`

const UpdateRequest = `
{
  "max_node_count": 20,
  "min_node_count": null,
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
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
	maxNodeCount = 20
	labels       = map[string]string{
		"gcloud_project_id":    "1",
		"gcloud_region_id":     "1",
		"gcloud_access_token":  "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4OTM3Nzc5LCJqdGkiOiI2NDk1MzFjMDNmYmU4NDczY2RiZGY2MjJkMWNmN2YzMCIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6MSwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxLCJ1c2VyX2lkIjoxLCJpc19hZG1pbiI6ZmFsc2V9.kdx_kpXF_Z_aCPDP5C-wIVLgv-mW9SSafJ_u6x7XO_k",
		"gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5OTExMzc3OSwianRpIjoiNjQ5NTMxYzAzZmJlODQ3M2NkYmRmNjIyZDFjZjdmMzAiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOjEsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MSwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.F_KTIGt1uvEaKnb8ZziI7Xca1o7Vcwj7752qbfb3Otg",
	}
	NodeGroupList1 = nodegroups.ClusterListNodeGroup{
		FlavorID:  "g1-standard-1-2",
		ImageID:   "fedora-coreos",
		NodeCount: 1,
		Name:      "default-master",
		UUID:      "3eda6b46-58d9-4abc-8a11-6045b791a35b",
		Role:      "master",
		Status:    "CREATE_IN_PROGRESS",
	}
	NodeGroupList2 = nodegroups.ClusterListNodeGroup{
		FlavorID:  "g1-standard-1-2",
		ImageID:   "fedora-coreos",
		NodeCount: 1,
		Name:      "default-worker",
		UUID:      "467fc654-2d1c-48e8-9d59-489e8fcf8c17",
		Role:      "worker",
		Status:    "CREATE_IN_PROGRESS",
	}
	NodeGroup1 = nodegroups.ClusterNodeGroup{
		ClusterID:            "c94f38cc-dc78-4715-8939-68de082bd5e3",
		ProjectID:            "46beed3938e6474390b530fefd6173d2",
		DockerVolumeSize:     10,
		Labels:               labels,
		NodeAddresses:        []net.IP{},
		MinNodeCount:         1,
		MaxNodeCount:         nil,
		IsDefault:            true,
		StackID:              "ccff14f2-4585-407d-aea5-8581aa9a292f",
		StatusReason:         nil,
		ClusterListNodeGroup: &NodeGroupList1,
	}
	UpdatedNodeGroup1 = nodegroups.ClusterNodeGroup{
		ClusterID:            "c94f38cc-dc78-4715-8939-68de082bd5e3",
		ProjectID:            "46beed3938e6474390b530fefd6173d2",
		DockerVolumeSize:     10,
		Labels:               labels,
		NodeAddresses:        []net.IP{},
		MinNodeCount:         1,
		MaxNodeCount:         &maxNodeCount,
		IsDefault:            true,
		StackID:              "ccff14f2-4585-407d-aea5-8581aa9a292f",
		StatusReason:         nil,
		ClusterListNodeGroup: &NodeGroupList1,
	}
	ExpectedClusterNodeGroupSlice = []nodegroups.ClusterListNodeGroup{NodeGroupList1, NodeGroupList2}
)

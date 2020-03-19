package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/volume/v1/volumes"
	"time"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "availability_zone": "nova",
      "created_at": "2019-05-29T05:32:41+0000",
      "volume_type": "standard",
      "id": "726ecfcc-7fd0-4e30-a86e-7892524aa483",
      "name": "123",
      "region": "Luxembourg 1",
      "status": "available",
      "updated_at": "2019-05-29T05:39:20+0000",
      "size": 2,
      "bootable": false,
      "project_id": 1,
      "region_id": 1,
      "attachments": [
        {
          "server_id": "8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
          "attachment_id": "f2ed59d9-8068-400c-be4b-c4501ef6f33c",
          "instance_name": "123",
          "attached_at": "2019-07-26T14:22:03+0000",
          "volume_id": "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
          "device": "/dev/vda"
        }
      ],
      "metadata": {
        "task_id": "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
        "attached_mode": "rw"
      },
      "creator_task_id": "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
      "volume_image_metadata": {
        "container_format": "bare",
        "min_ram": "0",
        "owner_specified.openstack.sha256": "87ddf8eea6504b5eb849e418a568c4985d3cea59b5a5d069e1dc644de676b4ec",
        "disk_format": "raw",
        "image_name": "cirros-gcloud",
        "image_id": "723037e2-ec6d-47eb-92de-6276c8907839",
        "owner_specified.openstack.object": "images/cirros-gcloud",
        "owner_specified.openstack.md5": "ba3cd24377dde5dfdd58728894004abb",
        "min_disk": "1",
        "checksum": "ba3cd24377dde5dfdd58728894004abb",
        "size": "46137344"
      }
    }
  ]
}
`

const GetResponse = `
{
  "availability_zone": "nova",
  "created_at": "2019-05-29T05:32:41+0000",
  "volume_type": "standard",
  "id": "726ecfcc-7fd0-4e30-a86e-7892524aa483",
  "name": "123",
  "region": "Luxembourg 1",
  "status": "available",
  "updated_at": "2019-05-29T05:39:20+0000",
  "size": 2,
  "bootable": false,
  "project_id": 1,
  "region_id": 1,
  "attachments": [
    {
      "server_id": "8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
      "attachment_id": "f2ed59d9-8068-400c-be4b-c4501ef6f33c",
      "instance_name": "123",
      "attached_at": "2019-07-26T14:22:03+0000",
      "volume_id": "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
      "device": "/dev/vda"
    }
  ],
  "metadata": {
    "task_id": "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
    "attached_mode": "rw"
  },
  "creator_task_id": "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
  "volume_image_metadata": {
    "container_format": "bare",
    "min_ram": "0",
    "owner_specified.openstack.sha256": "87ddf8eea6504b5eb849e418a568c4985d3cea59b5a5d069e1dc644de676b4ec",
    "disk_format": "raw",
    "image_name": "cirros-gcloud",
    "image_id": "723037e2-ec6d-47eb-92de-6276c8907839",
    "owner_specified.openstack.object": "images/cirros-gcloud",
    "owner_specified.openstack.md5": "ba3cd24377dde5dfdd58728894004abb",
    "min_disk": "1",
    "checksum": "ba3cd24377dde5dfdd58728894004abb",
    "size": "46137344"
  }
}
`

const CreateRequest = `
{
  "source": "new-volume",
  "type_name": "ssd_hiiops",
  "size": 10,
  "name": "TestVM5 Ubuntu volume",
  "instance_id_to_attach_to": "88f3e0bd-ca86-4cf7-be8b-dd2988e23c2d"
}	
`

const AttachDetachRequest = `
{
  "instance_id": "8dc30d49-bb34-4920-9bbd-03a2587ec0ad"
}	
`

const RetypeRequest = `
{
  "volume_type": "ssd_hiiops"
}	
`

const ExtendRequest = `
{
  "size": 16
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

const ExtendResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const createdTimeString = "2019-05-29T05:32:41+0000"
const updatedTimeString = "2019-05-29T05:39:20+0000"
const attachedTimeString = "2019-07-26T14:22:03+0000"

var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
var attachedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, attachedTimeString)
var attachedTime = gcorecloud.JSONRFC3339Z{Time: attachedTimeParsed}

var (
	Volume1 = volumes.Volume{
		AvailabilityZone: "nova",
		CreatedAt:        createdTime,
		UpdatedAt:        updatedTime,
		VolumeType:       "standard",
		ID:               "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		Name:             "123",
		RegionName:       "Luxembourg 1",
		Status:           "available",
		Size:             2,
		Bootable:         false,
		ProjectID:        1,
		RegionID:         1,
		Attachments: []volumes.Attachment{{
			ServerID:     "8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
			AttachmentID: "f2ed59d9-8068-400c-be4b-c4501ef6f33c",
			InstanceName: "123",
			AttachedAt:   attachedTime,
			VolumeID:     "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
			Device:       "/dev/vda",
		},
		},
		Metadata: volumes.Metadata{
			TaskID:       "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
			AttachedMode: "rw",
		},
		CreatorTaskID: "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
		VolumeImageMetadata: volumes.VolumeImageMetadata{
			ContainerFormat:               "bare",
			MinRAM:                        "0",
			OwnerSpecifiedOpenstackSHA256: "87ddf8eea6504b5eb849e418a568c4985d3cea59b5a5d069e1dc644de676b4ec",
			DiskFormat:                    "raw",
			ImageName:                     "cirros-gcloud",
			ImageID:                       "723037e2-ec6d-47eb-92de-6276c8907839",
			OwnerSpecifiedOpenstackObject: "images/cirros-gcloud",
			OwnerSpecifiedOpenstackMD5:    "ba3cd24377dde5dfdd58728894004abb",
			MinDisk:                       "1",
			Checksum:                      "ba3cd24377dde5dfdd58728894004abb",
			Size:                          "46137344",
		},
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedVolumeSlice = []volumes.Volume{Volume1}
)

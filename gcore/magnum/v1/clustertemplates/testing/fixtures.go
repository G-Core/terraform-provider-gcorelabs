package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/clustertemplates"
	"time"
)

const ListResponse = `
{
  "count": 1,	
  "results": [
    {
      "docker_volume_size": 5,
      "hidden": false,
      "fixed_subnet": null,
      "master_lb_enabled": true,
      "public": false,
      "apiserver_port": null,
      "docker_storage_driver": "overlay",
      "coe": "kubernetes",
      "network_driver": "calico",
      "cluster_distro": "fedora-coreos",
      "fixed_network": null,
      "links": [
        {
          "href": "http://10.100.178.165:9511/v1/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
          "rel": "self"
        },
        {
          "href": "http://10.100.178.165:9511/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
          "rel": "bookmark"
        }
      ],
      "https_proxy": null,
      "no_proxy": null,
      "registry_enabled": false,
      "labels": {},
      "server_type": "vm",
      "master_flavor_id": "g1-standard-1-2",
      "volume_driver": "cinder",
      "http_proxy": null,
      "name": "riswuwudlp",
      "external_network_id": "3a1d8f88-6020-4d11-9f9f-f815fe69ea40",
      "keypair_id": "keypair",
      "image_id": "fedora-coreos",
      "floating_ip_enabled": false,
      "dns_nameserver": "8.8.8.8",
      "user_id": "8ba64372-1585-4808-b422-7a7aab5f3197",
      "flavor_id": "g1-standard-1-2",
      "insecure_registry": null,
      "project_id": "eb83040c-35f5-4930-b5a5-b60c3fb016a6",
      "tls_disabled": false,
      "created_at": "2020-02-28T11:35:58+00:00",
      "uuid": "db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
      "updated_at": null
    }
  ]
}
`

const GetResponse = `
{
  "docker_volume_size": 5,
  "hidden": false,
  "fixed_subnet": null,
  "master_lb_enabled": true,
  "public": false,
  "apiserver_port": null,
  "docker_storage_driver": "overlay",
  "coe": "kubernetes",
  "network_driver": "calico",
  "cluster_distro": "fedora-coreos",
  "fixed_network": null,
  "links": [
    {
      "href": "http://10.100.178.165:9511/v1/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
      "rel": "self"
    },
    {
      "href": "http://10.100.178.165:9511/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
      "rel": "bookmark"
    }
  ],
  "https_proxy": null,
  "no_proxy": null,
  "registry_enabled": false,
  "labels": {},
  "server_type": "vm",
  "master_flavor_id": "g1-standard-1-2",
  "volume_driver": "cinder",
  "http_proxy": null,
  "name": "riswuwudlp",
  "external_network_id": "3a1d8f88-6020-4d11-9f9f-f815fe69ea40",
  "keypair_id": "keypair",
  "image_id": "fedora-coreos",
  "floating_ip_enabled": false,
  "dns_nameserver": "8.8.8.8",
  "user_id": "8ba64372-1585-4808-b422-7a7aab5f3197",
  "flavor_id": "g1-standard-1-2",
  "insecure_registry": null,
  "project_id": "eb83040c-35f5-4930-b5a5-b60c3fb016a6",
  "tls_disabled": false,
  "created_at": "2020-02-28T11:35:58+00:00",
  "uuid": "db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
  "updated_at": null
}
`

const CreateRequest = `
{
  "docker_volume_size": 5,
  "name": "riswuwudlp",
  "keypair_id": "keypair",
  "image_id": "fedora-coreos"
}
`

const CreateResponse = `
{
  "docker_volume_size": 5,
  "hidden": false,
  "fixed_subnet": null,
  "master_lb_enabled": true,
  "public": false,
  "apiserver_port": null,
  "docker_storage_driver": "overlay",
  "coe": "kubernetes",
  "network_driver": "calico",
  "cluster_distro": "fedora-coreos",
  "fixed_network": null,
  "links": [
    {
      "href": "http://10.100.178.165:9511/v1/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
      "rel": "self"
    },
    {
      "href": "http://10.100.178.165:9511/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
      "rel": "bookmark"
    }
  ],
  "https_proxy": null,
  "no_proxy": null,
  "registry_enabled": false,
  "labels": {},
  "server_type": "vm",
  "master_flavor_id": "g1-standard-1-2",
  "volume_driver": "cinder",
  "http_proxy": null,
  "name": "riswuwudlp",
  "external_network_id": "3a1d8f88-6020-4d11-9f9f-f815fe69ea40",
  "keypair_id": "keypair",
  "image_id": "fedora-coreos",
  "floating_ip_enabled": false,
  "dns_nameserver": "8.8.8.8",
  "user_id": "8ba64372-1585-4808-b422-7a7aab5f3197",
  "flavor_id": "g1-standard-1-2",
  "insecure_registry": null,
  "project_id": "eb83040c-35f5-4930-b5a5-b60c3fb016a6",
  "tls_disabled": false,
  "created_at": "2020-02-28T11:35:58+00:00",
  "uuid": "db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
  "updated_at": null
}
`

var createdTimeString = "2020-02-28T11:35:58+00:00"
var createdTime, _ = time.Parse(time.RFC3339, createdTimeString)

var (
	ClusterTemplate1 = clustertemplates.ClusterTemplate{
		Labels:              map[string]string{},
		FixedSubnet:         "",
		MasterFlavorID:      "g1-standard-1-2",
		FlavorID:            "g1-standard-1-2",
		NoProxy:             "",
		HTTPSProxy:          "",
		HTTPProxy:           "",
		TLSDisabled:         false,
		KeyPairID:           "keypair",
		Public:              false,
		DockerVolumeSize:    5,
		ServerType:          "vm",
		ExternalNetworkID:   "3a1d8f88-6020-4d11-9f9f-f815fe69ea40",
		ImageID:             "fedora-coreos",
		VolumeDriver:        "cinder",
		RegistryEnabled:     false,
		DockerStorageDriver: "overlay",
		Name:                "riswuwudlp",
		NetworkDriver:       "calico",
		FixedNetwork:        "",
		MasterLbEnabled:     true,
		DNSNameServer:       "8.8.8.8",
		FloatingIPEnabled:   false,
		Hidden:              false,
		UUID:                "db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
		CreatedAt:           createdTime,
		UpdatedAt:           nil,
		InsecureRegistry:    "",
		Links: []gcorecloud.Link{{
			Href: "http://10.100.178.165:9511/v1/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
			Rel:  "self",
		}, {
			Href: "http://10.100.178.165:9511/clustertemplates/db7a8580-cd3b-493f-a15a-8dbc9e30aa1e",
			Rel:  "bookmark",
		}},
	}

	ExpectedClusterTemplateSlice = []clustertemplates.ClusterTemplate{ClusterTemplate1}
)

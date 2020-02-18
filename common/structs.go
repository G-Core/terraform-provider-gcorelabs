package common

type Config struct {
	Jwt string
}


type InstanceId struct {
	InstanceId string `json:"instance_id"`
}

//K
type Keystone struct{
	Id int `json:"id"`
	Keystone_federated_domain_id string `json:"keystone_federated_domain_id"`
	State string `json:"state"`
	Admin_password string `json:"admin_password"`
	Created_on string `json:"created_on"`
	Url string `json:"url"`
}

//P
type Project struct {
	Id int `json:"id"`
	Client int `json:"client"`
	Name string `json:"name"`
	State string `json:"state"`
	Task_id int `string:"task_id"`
	Created_on string `string:"created_on"`
}

type Projects struct{
	Count int `json:"count"`
	Results []Project `json:"results"`
}


//R
type Region struct {
	Id int `json:"id"`
	Keystone_id int `json:"keystone_id"`
	State string `json:"state"`
	External_network_id string `json:"external_network_id"`
	Created_on string `json:"created_on"`
	Spice_proxy_url string `json:"spice_proxy_url"`
	Display_name string `json:"display_name"`
	Keystone_name string `json:"keystone_name"`
	Endpoint_type string `json:"endpoint_type"`
	Keystone_ Keystone `json:"keystone"`
}

type Regions struct{
	Count int `json:"count"`
	Results []Region `json:"results"`
}

//S
type Size struct{
	Size int `json:"size"`
}

//T
type TaskIds struct {
	Ids []string
}

type Type struct {
	Volume_type string `json:"volume_type"`
}

//V
type Volume struct{
	Size int `json:"size"`
	Source string `json:"source"`
	Name string	`json:"name"`
	Type_name string `json:"type_name"`
	Image_id string `json:"image_id"`
	Snapshot_id string `json:"snapshot_id"`
	Instance_id_to_attach_to string `json:"instance_id_to_attach_to"`
}

type VolumeAttachment struct{
	Server_id string `json:"server_id"`
	Instance_name string `json:"instance_name"`
	Attachment_id string `json:"attachment_id"`
	Volume_id string `json:"volume_id"`
	Device string `json:"device"`
	Attached_at string `json:"attached_at"`
}

type VolumeAttachments struct{
	Attachments []VolumeAttachment `json:"attachments"`
}

type VolumeIds struct {
	Ids []string
}


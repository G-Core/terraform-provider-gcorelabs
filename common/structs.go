package common

type Config struct {
	Jwt string
}


type InstanceId struct {
	InstanceId string `json:"instance_id"`
}

//P
type Project struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Projects struct{
	Count int `json:"count"`
	Results []Project `json:"results"`
}


//R
type Region struct {
	Id int `json:"id"`
	Keystone_name string `json:"keystone_name"`
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
	Ids []string `json:"tasks"`
}

type Type struct {
	Volume_type string `json:"volume_type"`
}

//V
type Volume struct{
	Size int `json:"size"`
	Source string `json:"source"`
	Name string	`json:"name"`
	Type_name string `json:"type_name,omitempty"`
	Image_id string `json:"image_id,omitempty"`
	Snapshot_id string `json:"snapshot_id,omitempty"`
	Instance_id_to_attach_to string `json:"instance_id_to_attach_to,omitempty"`
}

type OpenstackVolume struct{
	Size int `json:"size"`
	Type_name string `json:"type_name,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
}

type VolumeAttachment struct{
	Server_id string `json:"server_id"`
	Instance_name string `json:"instance_name"`
	Attachment_id string `json:"attachment_id"`
	Volume_id string `json:"volume_id"`
	Device string `json:"device"`
	Attached_at string `json:"attached_at"`
}

type VolumeIds struct {
	Volumes []string `json:"volumes"`
}


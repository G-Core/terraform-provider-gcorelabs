package common

type Config struct {
	Jwt string
}

type InstanceId struct {
	InstanceId string `json:"instance_id"`
}

//P
type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Projects struct {
	Count   int       `json:"count"`
	Results []Project `json:"results"`
}

//R
type Region struct {
	Id            int    `json:"id"`
	Keystone_name string `json:"keystone_name"`
}

type Regions struct {
	Count   int      `json:"count"`
	Results []Region `json:"results"`
}

//S
type Size struct {
	Size int `json:"size"`
}

//T
type TaskIds struct {
	Ids []string `json:"tasks"`
}

type Type struct {
	VolumeType string `json:"volume_type"`
}

//V
type Volume struct {
	Size       int    `json:"size"`
	Source     string `json:"source"`
	Name       string `json:"name"`
	TypeName   string `json:"type_name,omitempty"`
	ImageID    string `json:"image_id,omitempty"`
	SnapshotID string `json:"snapshot_id,omitempty"`
}

type OpenstackVolume struct {
	Size        int           `json:"size"`
	TypeName    string        `json:"volume_type,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
}

type VolumeIds struct {
	Volumes []string `json:"volumes"`
}

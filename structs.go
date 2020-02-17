package main

type Config struct {
	jwt string
}

type CreateVolumeBody struct{
	Size int `json:"size"`
	Source string `json:"source"`
	Name string	`json:"name"`
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

//T
type Task struct {
	Tasks []string
}

//V
type Volumes struct {
	Volumes []string
}


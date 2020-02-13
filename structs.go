package main

type Config struct {
	jwt string
}

type Volumes struct {
	Volumes []string
}

type Task struct {
	Tasks []string
}

type CreateVolumeBody struct{
	Size int `json:"size"`
	Source string `json:"source"`
	Name string	`json:"name"`
}
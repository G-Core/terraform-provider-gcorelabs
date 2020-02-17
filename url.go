package main

import "fmt"

func volume_url(project_id int, region_id int, volume_id string) string {
	return fmt.Sprintf("%svolumes/%d/%d/%s", HOST, project_id, region_id, volume_id)
}

func volumes_url(project_id int, region_id int) string {
	return fmt.Sprintf("%svolumes/%d/%d", HOST, project_id, region_id)
}



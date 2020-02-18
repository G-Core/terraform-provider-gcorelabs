package common

import "fmt"

func TaskUrl(task_id string) string {
	return fmt.Sprintf("%stasks/%s", HOST, task_id)
}

func VolumeUrl(project_id int, region_id int, volume_id string) string {
	return fmt.Sprintf("%svolumes/%d/%d/%s", HOST, project_id, region_id, volume_id)
}

func VolumesUrl(project_id int, region_id int) string {
	return fmt.Sprintf("%svolumes/%d/%d", HOST, project_id, region_id)
}



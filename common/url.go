package common

import "fmt"

func ExpandedObjectUrl(object_type string, project_id int, region_id int, volume_id string, addition string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s/%s", HOST, object_type, project_id, region_id, volume_id, addition)
}

func ObjectUrl(object_type string, project_id int, region_id int, volume_id string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s", HOST, object_type, project_id, region_id, volume_id)
}

// Get a list and post requests
func ObjectsUrl(object_type string, project_id int, region_id int) string {
	return fmt.Sprintf("%s%s/%d/%d", HOST, object_type, project_id, region_id)
}



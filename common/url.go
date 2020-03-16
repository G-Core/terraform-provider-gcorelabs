package common

import "fmt"

func ExpandedObjectUrl(object_type string, projectID int, regionID int, volumeID string, addition string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s/%s", HOST, object_type, projectID, regionID, volumeID, addition)
}

func ObjectUrl(object_type string, projectID int, regionID int, volumeID string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s", HOST, object_type, projectID, regionID, volumeID)
}

// Get a list and post requests
func ObjectsUrl(object_type string, projectID int, regionID int) string {
	return fmt.Sprintf("%s%s/%d/%d", HOST, object_type, projectID, regionID)
}

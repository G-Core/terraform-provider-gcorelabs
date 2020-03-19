package common

import "fmt"

func ExpandedObjectURL(objectType string, projectID int, regionID int, volumeID string, addition string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s/%s", HOST, objectType, projectID, regionID, volumeID, addition)
}

func ObjectURL(objectType string, projectID int, regionID int, volumeID string) string {
	return fmt.Sprintf("%s%s/%d/%d/%s", HOST, objectType, projectID, regionID, volumeID)
}

// ObjectsUrl creates an url for getting list requests and post requests
func ObjectsURL(objectType string, projectID int, regionID int) string {
	return fmt.Sprintf("%s%s/%d/%d", HOST, objectType, projectID, regionID)
}

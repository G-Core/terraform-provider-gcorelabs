package common

import "fmt"

// Example: "http://localhost:8888/v1/volumes/1/1/<uuid>/retype"
func ExpandedResourceV1URL(host string, objectType string, projectID int, regionID int, resourceID string, addition string) string {
	return fmt.Sprintf("%sv1/%s/%d/%d/%s/%s", host, objectType, projectID, regionID, resourceID, addition)
}

// Example: "http://localhost:8888/v1/volumes/1/1/<uuid>"
func ResourceV1URL(host string, objectType string, projectID int, regionID int, resourceID string) string {
	return fmt.Sprintf("%sv1/%s/%d/%d/%s", host, objectType, projectID, regionID, resourceID)
}

// ResourcesUrl creates an url for getting list requests and post requests
// Example: "http://localhost:8888/v1/volumes/1/1"
func ResourcesV1URL(host string, objectType string, projectID int, regionID int) string {
	return fmt.Sprintf("%sv1/%s/%d/%d", host, objectType, projectID, regionID)
}

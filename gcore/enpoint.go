package gcore

import (
	"gcloud/gcorecloud-go"
	"os"
	"strconv"
)

var nilEndpointOptions = gcorecloud.EndpointOpts{}

func EndpointOptionsFromEnv() (gcorecloud.EndpointOpts, error) {
	region := os.Getenv("GCLOUD_REGION")
	project := os.Getenv("GCLOUD_PROJECT")

	if region == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_REGION",
		}
		return nilEndpointOptions, err
	}

	regionInt, err := strconv.Atoi(region)
	if err != nil {
		return nilEndpointOptions, err
	}
	if project == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_PROJECT",
		}
		return nilEndpointOptions, err
	}

	projectInt, err := strconv.Atoi(project)
	if err != nil {
		return nilEndpointOptions, err
	}

	eo := gcorecloud.EndpointOpts{
		Region:  regionInt,
		Project: projectInt,
	}

	return eo, nil
}

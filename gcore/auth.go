package gcore

import (
	"gcloud/gcorecloud-go"
	"os"
	"strconv"
)

var nilAuthOptions = gcorecloud.AuthOptions{}
var nilTokenOptions = gcorecloud.TokenOptions{}

/*
AuthOptionsFromEnv fills out an identity.AuthOptions structure with the
settings found on environment variables.

The following variables provide sources of truth: GCLOUD_USERNAME, GCLOUD_PASSWORD, GCLOUD_AUTH_URL
	opts, err := gcore.AuthOptionsFromEnv()
	provider, err := gcore.AuthenticatedClient(opts)
*/
func AuthOptionsFromEnv() (gcorecloud.AuthOptions, error) {
	authURL := os.Getenv("GCLOUD_AUTH_URL")
	username := os.Getenv("GCLOUD_USERNAME")
	password := os.Getenv("GCLOUD_PASSWORD")

	if authURL == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_AUTH_URL",
		}
		return nilAuthOptions, err
	}

	if username == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_USERNAME",
		}
		return nilAuthOptions, err
	}

	if password == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_PASSWORD",
		}
		return nilAuthOptions, err
	}

	ao := gcorecloud.AuthOptions{
		IdentityEndpoint: authURL,
		Username:         username,
		Password:         password,
	}

	return ao, nil
}

func TokenOptionsFromEnv() (gcorecloud.TokenOptions, error) {

	apiURL := os.Getenv("GCLOUD_API_URL")
	accessToken := os.Getenv("GCLOUD_ACCESS_TOKEN")
	refreshToken := os.Getenv("GCLOUD_REFRESH_TOKEN")

	if apiURL == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_API_URL",
		}
		return nilTokenOptions, err
	}

	if accessToken == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_ACCESS_TOKEN",
		}
		return nilTokenOptions, err
	}

	if refreshToken == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_REFRESH_TOKEN",
		}
		return nilTokenOptions, err
	}

	to := gcorecloud.TokenOptions{
		IdentityEndpoint: apiURL,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AllowReauth:      true,
	}

	return to, nil
}

func NewGCloudTokenApiSettingsFromEnv() (*gcorecloud.GCloudTokenApiSettings, error) {
	apiURL := os.Getenv("GCLOUD_API_URL")
	apiVersion := os.Getenv("GCLOUD_API_VERSION")
	accessToken := os.Getenv("GCLOUD_ACCESS_TOKEN")
	refreshToken := os.Getenv("GCLOUD_REFRESH_TOKEN")
	region := os.Getenv("GCLOUD_REGION")
	project := os.Getenv("GCLOUD_PROJECT")

	var projectInt, regionInt int
	var err error
	var version = "v1"

	if project != "" {
		projectInt, err = strconv.Atoi(project)
		if err != nil {
			return nil, err
		}
	}

	if region != "" {
		regionInt, err = strconv.Atoi(region)
		if err != nil {
			return nil, err
		}
	}

	if apiVersion != "" {
		version = apiVersion
	}

	return &gcorecloud.GCloudTokenApiSettings{
		Version:          version,
		IdentityEndpoint: apiURL,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		Region:           regionInt,
		Project:          projectInt,
		AllowReauth:      true,
	}, nil
}

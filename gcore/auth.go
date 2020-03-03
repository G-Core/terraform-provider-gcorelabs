package gcore

import (
	"gcloud/gcorecloud-go"
	"os"
)

var nilOptions = gcorecloud.AuthOptions{}
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
		return nilOptions, err
	}

	if username == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_USERNAME",
		}
		return nilOptions, err
	}

	if password == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_PASSWORD",
		}
		return nilOptions, err
	}

	ao := gcorecloud.AuthOptions{
		IdentityEndpoint: authURL,
		Username: username,
		Password: password,
	}

	return ao, nil
}

func TokenOptionsFromEnv() (gcorecloud.TokenOptions, error) {

	refreshURL := os.Getenv("GCLOUD_REFRESH_URL")
	accessToken := os.Getenv("GCLOUD_ACCESS_TOKEN")
	refreshToken := os.Getenv("GCLOUD_REFRESH_TOKEN")

	if refreshURL == "" {
		err := gcorecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "GCLOUD_REFRESH_URL",
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
		IdentityEndpoint: refreshURL,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return to, nil
}

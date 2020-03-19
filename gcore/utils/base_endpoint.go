package utils

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"net/url"
	"regexp"
	"strings"
)

// BaseEndpoint will return a URL without the /vX.Y
// portion of the URL.
func BaseEndpoint(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	u.RawQuery, u.Fragment = "", ""

	path := u.Path
	versionRe := regexp.MustCompile("v[0-9.]+/?")

	if version := versionRe.FindString(path); version != "" {
		versionIndex := strings.Index(path, version)
		u.Path = path[:versionIndex]
	}

	return u.String(), nil
}

// BaseVersionEndpoint will return a URL with the /vX.Y starting portion of the URL.
func BaseVersionEndpoint(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	u.RawQuery, u.Fragment = "", ""
	u.Path = gcorecloud.NormalizeURL(strings.Join(strings.Split(u.Path, "/")[:2], "/"))
	return u.String(), nil
}

func BaseRootEndpoint(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	u.RawQuery, u.Fragment, u.Path = "", "", ""

	return u.String(), nil
}

// NormalizeURLPath removes duplicated slashes
func NormalizeURLPath(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	path := u.Path
	r := regexp.MustCompile(`//+`)
	u.Path = r.ReplaceAllLiteralString(path, "/")
	return gcorecloud.NormalizeURL(u.String()), nil
}

package client

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

// Fake token to use.
const TokenID = "cbc36478b0bd8e67e89469c7749d4127"      // nolint
const AccessToken = "cbc36478b0bd8e67e89469c7749d4127"  // nolint
const RefreshToken = "tbc36478b0bd8e67e89469c7749d4127" // nolint
const Username = "username"                             // nolint
const Password = "password"                             // nolint
const RegionID = 1                                      // nolint
const ProjectID = 1                                     // nolint

// ServiceClient returns a generic service client for use in tests.
func ServiceClient() *gcorecloud.ServiceClient {
	return &gcorecloud.ServiceClient{
		ProviderClient: &gcorecloud.ProviderClient{
			AccessTokenID:  AccessToken,
			RefreshTokenID: RefreshToken,
		},
		Endpoint: testhelper.Endpoint(),
	}
}

func ServiceTokenClient(name string, version string) *gcorecloud.ServiceClient {
	options := gcorecloud.TokenOptions{
		APIURL:       testhelper.GCloudRefreshTokenIdentifyEndpoint(),
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		AllowReauth:  true,
	}
	endpointOpts := gcorecloud.EndpointOpts{
		Name:    name,
		Region:  RegionID,
		Project: ProjectID,
		Version: version,
	}
	client, err := gcore.TokenClientService(options, endpointOpts)
	if err != nil {
		panic(err)
	}
	return client
}

func ServiceAuthClient(name string, version string) *gcorecloud.ServiceClient {
	options := gcorecloud.AuthOptions{
		APIURL:      testhelper.GCoreIdentifyEndpoint(),
		AuthURL:     testhelper.GCoreRefreshTokenIdentifyEndpoint(),
		Username:    Username,
		Password:    Password,
		AllowReauth: true,
	}
	endpointOpts := gcorecloud.EndpointOpts{
		Name:    name,
		Region:  RegionID,
		Project: ProjectID,
		Version: version,
	}
	client, err := gcore.AuthClientService(options, endpointOpts)
	if err != nil {
		panic(err)
	}
	return client
}

type AuthResultTest struct {
	accessToken  string
	refreshToken string
}

func (ar AuthResultTest) ExtractAccessToken() (string, error) {
	return ar.accessToken, nil
}

func (ar AuthResultTest) ExtractRefreshToken() (string, error) {
	return ar.accessToken, nil
}

func (ar AuthResultTest) ExtractTokensPair() (string, string, error) {
	return ar.accessToken, ar.refreshToken, nil
}

func NewAuthResultTest(accessToken string, refreshToken string) AuthResultTest {
	return AuthResultTest{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

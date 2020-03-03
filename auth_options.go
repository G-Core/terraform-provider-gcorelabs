package gcorecloud

import "fmt"

type AuthOptionsBuilder interface {
	ToMap() map[string]interface{}
}

type TokenOptionsBuilder interface {
	ToMap() map[string]interface{}
}

type AuthOptions struct {
	IdentityEndpoint     string `json:"-"`
	RefreshTokenEndpoint string `json:"-"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	AllowReauth          bool   `json:"-"`
}

func (ao AuthOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"username": ao.Username,
		"password": ao.Password,
	}
}

type TokenOptions struct {
	IdentityEndpoint string `json:"-"`
	AccessToken      string `json:"access,omitempty"`
	RefreshToken     string `json:"refresh,omitempty"`
	AllowReauth      bool   `json:"-"`
}

func (to TokenOptions) ExtractAccessToken() (string, error) {
	return to.AccessToken, nil
}
func (to TokenOptions) ExtractRefreshToken() (string, error) {
	return to.RefreshToken, nil
}
func (to TokenOptions) ExtractTokensPair() (string, string, error) {
	return to.AccessToken, to.RefreshToken, nil
}

func (to TokenOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"token": to.RefreshToken,
	}
}

type GCloudTokenApiSettings struct {
	IdentityEndpoint string `json:"url,omitempty"`
	AccessToken      string `json:"access,omitempty"`
	RefreshToken     string `json:"refresh,omitempty"`
	AllowReauth      bool   `json:"-"`
	Type             string `json:"type,omitempty"`
	Name             string `json:"name,omitempty"`
	Region           int    `json:"region,omitempty"`
	Project          int    `json:"project,omitempty"`
	Version          string `json:"version,omitempty"`
}

func (gs GCloudTokenApiSettings) ToTokenOptions() TokenOptions {
	return TokenOptions{
		IdentityEndpoint: gs.IdentityEndpoint,
		AccessToken:      gs.AccessToken,
		RefreshToken:     gs.RefreshToken,
		AllowReauth:      gs.AllowReauth,
	}
}

func (gs GCloudTokenApiSettings) ToEndpointOptions() EndpointOpts {
	return EndpointOpts{
		Region:  gs.Region,
		Project: gs.Project,
		Version: gs.Version,
		Name:    gs.Name,
		Type:    gs.Type,
	}
}

func (gs GCloudTokenApiSettings) Validate() error {
	if gs.AccessToken == "" {
		return fmt.Errorf("access token required")
	}
	if gs.RefreshToken == "" {
		return fmt.Errorf("refresh token required")
	}
	if gs.IdentityEndpoint == "" {
		return fmt.Errorf("api url required. IdentityEndpoint")
	}
	if gs.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

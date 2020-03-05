package gcorecloud

import "fmt"

type AuthOptionsBuilder interface {
	ToMap() map[string]interface{}
}

type TokenOptionsBuilder interface {
	ToMap() map[string]interface{}
}

type AuthOptions struct {
	ApiURL      string `json:"-"`
	AuthURL     string `json:"-"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AllowReauth bool   `json:"-"`
}

func (ao AuthOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"username": ao.Username,
		"password": ao.Password,
	}
}

type TokenOptions struct {
	ApiURL       string `json:"-"`
	AccessToken  string `json:"access,omitempty"`
	RefreshToken string `json:"refresh,omitempty"`
	AllowReauth  bool   `json:"-"`
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
	ApiURL       string `json:"url,omitempty"`
	AccessToken  string `json:"access,omitempty"`
	RefreshToken string `json:"refresh,omitempty"`
	AllowReauth  bool   `json:"-"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	Region       int    `json:"region,omitempty"`
	Project      int    `json:"project,omitempty"`
	Version      string `json:"version,omitempty"`
	Debug        bool   `json:"debug,omitempty"`
}

func (gs GCloudTokenApiSettings) ToTokenOptions() TokenOptions {
	return TokenOptions{
		ApiURL:       gs.ApiURL,
		AccessToken:  gs.AccessToken,
		RefreshToken: gs.RefreshToken,
		AllowReauth:  gs.AllowReauth,
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
	if gs.ApiURL == "" {
		return fmt.Errorf("api-url endpoint required")
	}
	if gs.AccessToken == "" {
		return fmt.Errorf("access token required")
	}
	if gs.RefreshToken == "" {
		return fmt.Errorf("refresh token required")
	}
	if gs.ApiURL == "" {
		return fmt.Errorf("api url required. ApiURL")
	}
	if gs.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

type GCloudPasswordApiSettings struct {
	ApiURL      string `json:"api-url,omitempty"`
	AuthURL     string `json:"auth-url,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AllowReauth bool   `json:"-"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Region      int    `json:"region,omitempty"`
	Project     int    `json:"project,omitempty"`
	Version     string `json:"version,omitempty"`
	Debug       bool   `json:"debug,omitempty"`
}

func (gs GCloudPasswordApiSettings) ToAuthOptions() AuthOptions {
	return AuthOptions{
		ApiURL:      gs.ApiURL,
		AuthURL:     gs.AuthURL,
		Username:    gs.Username,
		Password:    gs.Password,
		AllowReauth: gs.AllowReauth,
	}
}

func (gs GCloudPasswordApiSettings) ToEndpointOptions() EndpointOpts {
	return EndpointOpts{
		Region:  gs.Region,
		Project: gs.Project,
		Version: gs.Version,
		Name:    gs.Name,
		Type:    gs.Type,
	}
}

func (gs GCloudPasswordApiSettings) Validate() error {
	if gs.AuthURL == "" {
		return fmt.Errorf("auth-url endpoint required")
	}
	if gs.ApiURL == "" {
		return fmt.Errorf("api-url endpoint required")
	}
	if gs.Username == "" {
		return fmt.Errorf("username required")
	}
	if gs.Password == "" {
		return fmt.Errorf("password required")
	}
	if gs.ApiURL == "" {
		return fmt.Errorf("api url required. ApiURL")
	}
	if gs.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

package gcorecloud

import "fmt"

// AuthOptionsBuilder build auth options to map
type AuthOptionsBuilder interface {
	ToMap() map[string]interface{}
}

// TokenOptionsBuilder build token options to map
type TokenOptionsBuilder interface {
	ToMap() map[string]interface{}
}

// AuthOptions gcore API
type AuthOptions struct {
	APIURL      string `json:"-"`
	AuthURL     string `json:"-"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AllowReauth bool   `json:"-"`
}

// ToMap implements AuthOptionsBuilder
func (ao AuthOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"username": ao.Username,
		"password": ao.Password,
	}
}

// TokenOptions gcore API
type TokenOptions struct {
	APIURL       string `json:"-"`
	AccessToken  string `json:"access,omitempty"`
	RefreshToken string `json:"refresh,omitempty"`
	AllowReauth  bool   `json:"-"`
}

// ExtractAccessToken implements AuthResult
func (to TokenOptions) ExtractAccessToken() (string, error) {
	return to.AccessToken, nil
}

// ExtractRefreshToken implements AuthResult
func (to TokenOptions) ExtractRefreshToken() (string, error) {
	return to.RefreshToken, nil
}

// ExtractTokensPair implements AuthResult
func (to TokenOptions) ExtractTokensPair() (string, string, error) {
	return to.AccessToken, to.RefreshToken, nil
}

// ToMap implements TokenOptionsBuilder
func (to TokenOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"token": to.RefreshToken,
	}
}

// TokenClientSettings interface
type TokenClientSettings interface {
	ToTokenOptions() TokenOptions
	ToEndpointOptions() EndpointOpts
	Validate() error
}

// AuthClientSettings interface
type AuthClientSettings interface {
	ToAuthOptions() AuthOptions
	ToEndpointOptions() EndpointOpts
	Validate() error
}

// TokenAPISettings - settings for token client building
type TokenAPISettings struct {
	APIURL       string `json:"url,omitempty"`
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

// ToTokenOptions implements TokenClientSettings interface
func (gs TokenAPISettings) ToTokenOptions() TokenOptions {
	return TokenOptions{
		APIURL:       gs.APIURL,
		AccessToken:  gs.AccessToken,
		RefreshToken: gs.RefreshToken,
		AllowReauth:  gs.AllowReauth,
	}
}

// ToEndpointOptions implements TokenClientSettings interface
func (gs TokenAPISettings) ToEndpointOptions() EndpointOpts {
	return EndpointOpts{
		Region:  gs.Region,
		Project: gs.Project,
		Version: gs.Version,
		Name:    gs.Name,
		Type:    gs.Type,
	}
}

// Validate implements TokenClientSettings interface
func (gs TokenAPISettings) Validate() error {
	if gs.APIURL == "" {
		return fmt.Errorf("api-url endpoint required")
	}
	if gs.AccessToken == "" {
		return fmt.Errorf("access token required")
	}
	if gs.RefreshToken == "" {
		return fmt.Errorf("refresh token required")
	}
	if gs.APIURL == "" {
		return fmt.Errorf("api url required. APIURL")
	}
	if gs.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

// PasswordAPISettings - settings for password client building
type PasswordAPISettings struct {
	APIURL      string `json:"api-url,omitempty"`
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

// ToAuthOptions implements AuthClientSettings interface
func (gs PasswordAPISettings) ToAuthOptions() AuthOptions {
	return AuthOptions{
		APIURL:      gs.APIURL,
		AuthURL:     gs.AuthURL,
		Username:    gs.Username,
		Password:    gs.Password,
		AllowReauth: gs.AllowReauth,
	}
}

// ToEndpointOptions implements AuthClientSettings interface
func (gs PasswordAPISettings) ToEndpointOptions() EndpointOpts {
	return EndpointOpts{
		Region:  gs.Region,
		Project: gs.Project,
		Version: gs.Version,
		Name:    gs.Name,
		Type:    gs.Type,
	}
}

// Validate implements AuthClientSettings interface
func (gs PasswordAPISettings) Validate() error {
	if gs.AuthURL == "" {
		return fmt.Errorf("auth-url endpoint required")
	}
	if gs.APIURL == "" {
		return fmt.Errorf("api-url endpoint required")
	}
	if gs.Username == "" {
		return fmt.Errorf("username required")
	}
	if gs.Password == "" {
		return fmt.Errorf("password required")
	}
	if gs.APIURL == "" {
		return fmt.Errorf("api url required. APIURL")
	}
	if gs.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

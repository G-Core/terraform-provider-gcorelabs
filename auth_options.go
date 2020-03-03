package gcorecloud

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

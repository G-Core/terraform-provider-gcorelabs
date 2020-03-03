package tokens

import (
	"gcloud/gcorecloud-go"
)

// commonResult is the response from a request. A commonResult has various
// methods which can be used to extract different details about the result.
type commonResult struct {
	gcorecloud.Result
}

// ExtractToken interprets a commonResult as a Token.
func (r commonResult) ExtractTokens() (*Token, error) {
	var s Token
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return &s, err
}

func (r TokenResult) ExtractAccessToken() (string, error) {
	t, err := r.ExtractTokens()
	if err != nil {
		return "", err
	}
	return t.Access, nil
}

func (r TokenResult) ExtractRefreshToken() (string, error) {
	t, err := r.ExtractTokens()
	if err != nil {
		return "", err
	}
	return t.Refresh, nil
}

func (r TokenResult) ExtractTokensPair() (string, string, error) {
	t, err := r.ExtractTokens()
	if err != nil {
		return "", "", err
	}
	return t.Access, t.Refresh, nil
}

// TokenResult is the response from a Create request. Use ExtractToken() to interpret it as a Token
type TokenResult struct {
	commonResult
}

type Token struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.ExtractIntoStructPtr(v, "")
}

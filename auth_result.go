package gcorecloud

type AuthResult interface {
	ExtractAccessToken() (string, error)
	ExtractRefreshToken() (string, error)
	ExtractTokensPair() (string, string, error)
}

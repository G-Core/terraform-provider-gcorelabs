package gcore

import (
	"reflect"

	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/identity/tokens"
	"gcloud/gcorecloud-go/gcore/utils"
)

/*
NewClient prepares an unauthenticated ProviderClient instance.
Most users will probably prefer using the AuthenticatedClient function instead.

This is useful if you wish to explicitly control the version of the identity
service that's used for authentication explicitly, for example.

A basic example of using this would be:

	ao, err := gcore.AuthOptionsFromEnv()
	provider, err := gcore.NewClient(ao.APIURL)
	client, err := gcore.NewIdentity(provider, gcorecloud.EndpointOpts{})
*/
func NewClient(endpoint string) (*gcorecloud.ProviderClient, error) {
	base, err := utils.BaseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpoint = gcorecloud.NormalizeURL(endpoint)
	base = gcorecloud.NormalizeURL(base)

	p := gcorecloud.NewProviderClient()
	p.IdentityBase = base
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

func NewGCoreClient(endpoint string) (*gcorecloud.ProviderClient, error) {
	base, err := utils.BaseRootEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpoint = gcorecloud.NormalizeURL(endpoint)
	base = gcorecloud.NormalizeURL(base)

	p := gcorecloud.NewProviderClient()
	p.IdentityBase = base
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

/*
AuthenticatedClient logs in to an GCore cloud found at the identity endpoint
specified by the options, acquires a token, and returns a Provider Client
instance that's ready to operate.

Example:

	ao, err := gcore.AuthOptionsFromEnv()
	provider, err := gcore.AuthenticatedClient(ao)
	client, err := gcore.NewMagnumV1(client, gcorecloud.EndpointOpts{})
*/
func AuthenticatedClient(options gcorecloud.AuthOptions) (*gcorecloud.ProviderClient, error) {
	client, err := NewGCoreClient(options.APIURL)
	if err != nil {
		return nil, err
	}
	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TokenClient(options gcorecloud.TokenOptions) (*gcorecloud.ProviderClient, error) {
	client, err := NewGCoreClient(options.APIURL)
	if err != nil {
		return nil, err
	}
	err = client.SetTokensAndAuthResult(options)
	if err != nil {
		return nil, err
	}
	setGCloudReauth(client, "", options, gcorecloud.EndpointOpts{})
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func Authenticate(client *gcorecloud.ProviderClient, options gcorecloud.AuthOptions) error {
	return auth(client, "", options, gcorecloud.EndpointOpts{})
}

func auth(client *gcorecloud.ProviderClient, endpoint string, options gcorecloud.AuthOptions, eo gcorecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		identityClient.Endpoint = endpoint
	}

	result := tokens.Create(identityClient, options)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		err = tac.SetTokensAndAuthResult(nil)
		if err != nil {
			return err
		}
		tro := client.ToTokenOptions()
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refresh(&tac, endpoint, tro, tao, eo)
			if err != nil {
				errAuth := auth(&tac, endpoint, tao, eo)
				if errAuth != nil {
					return errAuth
				}
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func refresh(client *gcorecloud.ProviderClient, endpoint string, tokenOptions gcorecloud.TokenOptions, authOptions gcorecloud.AuthOptions, eo gcorecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		identityClient.Endpoint = endpoint
	}

	result := tokens.Refresh(identityClient, tokenOptions)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if tokenOptions.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tro := tokenOptions
		tro.AllowReauth = false
		tao := authOptions
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refresh(&tac, endpoint, tro, tao, eo)
			if err != nil {
				errAuth := auth(&tac, endpoint, tao, eo)
				if errAuth != nil {
					return errAuth
				}
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func refreshGCloud(client *gcorecloud.ProviderClient, endpoint string, options gcorecloud.TokenOptions, eo gcorecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		identityClient.Endpoint = endpoint
	}

	result := tokens.RefreshGCloud(identityClient, options)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshGCloud(&tac, endpoint, tao, eo)
			if err != nil {
				return err
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func setGCloudReauth(client *gcorecloud.ProviderClient, endpoint string, options gcorecloud.TokenOptions, eo gcorecloud.EndpointOpts) {

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshGCloud(&tac, endpoint, tao, eo)
			if err != nil {
				return err
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}
}

// NewIdentity creates a ServiceClient that may be used to interact with the gcore identity auth service.
func NewIdentity(client *gcorecloud.ProviderClient, eo gcorecloud.EndpointOpts) (*gcorecloud.ServiceClient, error) {
	endpoint := client.IdentityBase
	clientType := "auth"
	var err error
	if !reflect.DeepEqual(eo, gcorecloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &gcorecloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

func initClientOpts(client *gcorecloud.ProviderClient, eo gcorecloud.EndpointOpts, clientType string) (*gcorecloud.ServiceClient, error) {
	sc := new(gcorecloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	url, err = utils.NormalizeURLPath(url)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	endpoint, err := utils.BaseVersionEndpoint(url)
	if err != nil {
		return sc, err
	}
	sc.Endpoint = endpoint
	sc.ResourceBase = url
	sc.Type = clientType
	return sc, nil
}

func TokenClientService(options gcorecloud.TokenOptions, eo gcorecloud.EndpointOpts) (*gcorecloud.ServiceClient, error) {
	provider, err := TokenClient(options)
	if err != nil {
		return nil, err
	}
	provider.EndpointLocator = gcorecloud.DefaultEndpointLocator(provider.IdentityBase)
	client, err := initClientOpts(provider, eo, eo.Type)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func AuthClientService(options gcorecloud.AuthOptions, eo gcorecloud.EndpointOpts) (*gcorecloud.ServiceClient, error) {
	provider, err := AuthenticatedClient(options)
	if err != nil {
		return nil, err
	}
	provider.EndpointLocator = gcorecloud.DefaultEndpointLocator(provider.IdentityBase)
	client, err := initClientOpts(provider, eo, eo.Type)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TaskTokenClient(options gcorecloud.TokenOptions) (*gcorecloud.ServiceClient, error) {
	eo := gcorecloud.EndpointOpts{
		Name:    "tasks",
		Version: "v1",
	}
	client, err := TokenClientService(options, eo)
	if err != nil {
		return nil, err
	}
	return client, nil
}

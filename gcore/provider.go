package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	dnssdk "github.com/G-Core/g-dns-sdk-go"
	storageSDK "github.com/G-Core/gcorelabs-storage-sdk-go"
	gcdn "github.com/G-Core/gcorelabscdn-go"
	gcdnProvider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ProviderOptPermanentToken   = "permanent_api_token"
	ProviderOptSkipCredsAuthErr = "ignore_creds_auth_error"

	lifecyclePolicyResource = "gcore_lifecyclepolicy"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCORE_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PASSWORD", ""),
			},
			ProviderOptPermanentToken: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "A permanent API-token. Implemented for Storage Terraform Resource only. https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PERMANENT_TOKEN", ""),
			},
			ProviderOptSkipCredsAuthErr: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Should be set to true when you are gonna to use storage resource with permanent API-token only.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PERMANENT_TOKEN", ""),
			},
			"gcore_platform": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Platform ulr is used for generate jwt",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PLATFORM", "https://api.gcdn.co"),
			},
			"gcore_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region API",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_API", "https://api.cloud.gcorelabs.com"),
			},
			"gcore_cdn_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CDN API",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_CDN_API", "https://api.gcdn.co/cdn"),
			},
			"gcore_storage_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage API",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_STORAGE_API", "https://storage.gcorelabs.com/api"),
			},
			"gcore_dns_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS API",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_DNS_API", "https://dnsapi.gcorelabs.com"),
			},
			"gcore_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client id",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_CLIENT_ID", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gcore_volume":           resourceVolume(),
			"gcore_network":          resourceNetwork(),
			"gcore_subnet":           resourceSubnet(),
			"gcore_router":           resourceRouter(),
			"gcore_instance":         resourceInstance(),
			"gcore_keypair":          resourceKeypair(),
			"gcore_reservedfixedip":  resourceReservedFixedIP(),
			"gcore_floatingip":       resourceFloatingIP(),
			"gcore_loadbalancer":     resourceLoadBalancer(),
			"gcore_lblistener":       resourceLbListener(),
			"gcore_lbpool":           resourceLBPool(),
			"gcore_lbmember":         resourceLBMember(),
			"gcore_securitygroup":    resourceSecurityGroup(),
			"gcore_baremetal":        resourceBmInstance(),
			"gcore_snapshot":         resourceSnapshot(),
			"gcore_servergroup":      resourceServerGroup(),
			"gcore_k8s":              resourceK8s(),
			"gcore_k8s_pool":         resourceK8sPool(),
			"gcore_storage_s3":       resourceStorageS3(),
			DNSZoneResource:          resourceDNSZone(),
			DNSZoneRecordResource:    resourceDNSZoneRecord(),
			"gcore_storage_sftp":     resourceStorageSFTP(),
			"gcore_storage_sftp_key": resourceStorageSFTPKey(),
			"gcore_cdn_resource":     resourceCDNResource(),
			"gcore_cdn_origingroup":  resourceCDNOriginGroup(),
			"gcore_cdn_rule":         resourceCDNRule(),
			"gcore_cdn_sslcert":      resourceCDNCert(),
			lifecyclePolicyResource:  resourceLifecyclePolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gcore_project":          dataSourceProject(),
			"gcore_region":           dataSourceRegion(),
			"gcore_securitygroup":    dataSourceSecurityGroup(),
			"gcore_image":            dataSourceImage(),
			"gcore_volume":           dataSourceVolume(),
			"gcore_network":          dataSourceNetwork(),
			"gcore_subnet":           dataSourceSubnet(),
			"gcore_router":           dataSourceRouter(),
			"gcore_loadbalancer":     dataSourceLoadBalancer(),
			"gcore_lblistener":       dataSourceLBListener(),
			"gcore_lbpool":           dataSourceLBPool(),
			"gcore_instance":         dataSourceInstance(),
			"gcore_floatingip":       dataSourceFloatingIP(),
			"gcore_storage_s3":       dataSourceStorageS3(),
			"gcore_storage_sftp":     dataSourceStorageSFTP(),
			"gcore_storage_sftp_key": dataSourceStorageSFTPKey(),
			"gcore_reservedfixedip":  dataSourceReservedFixedIP(),
			"gcore_servergroup":      dataSourceServerGroup(),
			"gcore_k8s":              dataSourceK8s(),
			"gcore_k8s_pool":         dataSourceK8sPool(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("user_name").(string)
	password := d.Get("password").(string)
	permanentToken := d.Get(ProviderOptPermanentToken).(string)
	api := d.Get("gcore_api").(string)
	cdnAPI := d.Get("gcore_cdn_api").(string)
	storageAPI := d.Get("gcore_storage_api").(string)
	dnsAPI := d.Get("gcore_dns_api").(string)
	platform := d.Get("gcore_platform").(string)
	clientID := d.Get("gcore_client_id").(string)

	var diags diag.Diagnostics

	provider, err := gc.AuthenticatedClient(gcorecloud.AuthOptions{
		APIURL:      api,
		AuthURL:     platform,
		Username:    username,
		Password:    password,
		AllowReauth: true,
		ClientID:    clientID,
	})

	skipAuthErr, ok := d.GetOk(ProviderOptSkipCredsAuthErr)
	if err != nil && !(ok == true || skipAuthErr.(bool) == true) {
		return nil, diag.FromErr(fmt.Errorf("init auth client: %w", err))

	}
	if err != nil {
		provider = &gcorecloud.ProviderClient{}
		log.Printf("[WARN] init auth client: %s\n", err)
	}

	cdnProvider := gcdnProvider.NewClient(cdnAPI, gcdnProvider.WithSignerFunc(func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+provider.AccessToken())
		return nil
	}))
	cdnService := gcdn.NewService(cdnProvider)

	config := Config{
		Provider:  provider,
		CDNClient: cdnService,
	}

	if storageAPI != "" {
		stHost, stPath, err := ExtractHostAndPath(storageAPI)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("storage api url: %w", err))
		}
		config.StorageClient = storageSDK.NewSDK(
			stHost,
			stPath,
			storageSDK.WithBearerAuth(provider.AccessToken),
			storageSDK.WithPermanentTokenAuth(func() string { return permanentToken }),
		)
	}
	if dnsAPI != "" {
		baseUrl, err := url.Parse(dnsAPI)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("dns api url: %w", err))
		}
		authorizer := dnssdk.BearerAuth(provider.AccessToken())
		if permanentToken != "" {
			authorizer = dnssdk.PermanentAPIKeyAuth(permanentToken)
		}
		config.DNSClient = dnssdk.NewClient(authorizer, func(client *dnssdk.Client) {
			client.BaseURL = baseUrl
		})
	}

	return &config, diags
}

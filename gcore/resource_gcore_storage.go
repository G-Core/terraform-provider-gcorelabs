package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/G-Core/gcorelabs-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageSchemaGenerateSftpPassword = "generate_sftp_password"
	StorageSchemaGenerateS3AccessKey  = "generate_s3_access_key"
	StorageSchemaGenerateS3SecretKey  = "generate_s3_secret_key"
	StorageSchemaLocation             = "location"
	StorageSchemaName                 = "name"
	StorageSchemaType                 = "type"
	StorageSchemaId                   = "storage_id"
	StorageSchemaClientId             = "client_id"
	StorageSchemaSftpPassword         = "sftp_password"
	StorageSchemaKeyId                = "ssh_key_id"
	StorageSchemaExpires              = "expires"
	StorageSchemaServerAlias          = "server_alias"
)

func resourceStorage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "An id of new storage resource.",
			},
			StorageSchemaClientId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "An client id of new storage resource.",
			},
			StorageSchemaName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A name of new storage resource.",
			},
			StorageSchemaServerAlias: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An alias of storage resource.",
			},
			StorageSchemaExpires: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A expires date of storage resource.",
			},
			StorageSchemaType: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					val := v.(string)
					allowed := []string{"sftp", "s3"}
					for _, el := range allowed {
						if el == val {
							return nil
						}
					}
					return diag.Errorf(`must be one of %+v`, allowed)
				},
				Description: "A type of new storage resource. One of (sftp, s3)",
			},
			StorageSchemaLocation: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					val := v.(string)
					allowed := []string{"s-ed1", "s-darz1", "s-ws1", "ams", "sin", "fra", "mia"}
					for _, el := range allowed {
						if el == val {
							return nil
						}
					}
					return diag.Errorf(`must be one of %+v`, allowed)
				},
				Description: "A location of new storage resource. One of (s-ed1, s-darz1, s-ws1, ams, sin, fra, mia)",
			},
			StorageSchemaSftpPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A sftp password for new storage resource.",
			},
			StorageSchemaGenerateS3AccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 access key for new storage resource.",
			},
			StorageSchemaGenerateS3SecretKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 secret key for new storage resource.",
			},
			StorageSchemaGenerateSftpPassword: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "An auto generated sftp password for new storage resource.",
			},
			StorageSchemaKeyId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An ssh key id to link with new sftp storage resource only. https://storage.gcorelabs.com/ssh-key/list",
			},
		},
		CreateContext: resourceStorageCreate,
		ReadContext:   resourceStorageRead,
		UpdateContext: resourceStorageUpdate,
		DeleteContext: resourceStorageDelete,
		Description:   "Represent storage resource. https://storage.gcorelabs.com/storage/list",
	}
}

func resourceStorageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (dErr diag.Diagnostics) {
	id := new(int)
	log.Println("[DEBUG] Start Storage Resource creating")
	defer log.Printf("[DEBUG] Finish Storage Resource creating (id=%d)\n", *id)
	config := m.(*Config)
	client := config.StorageClient

	opts := make([]func(opt *storage.StorageCreateHTTPParams), 0)
	opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Context = ctx })
	if d.Get(StorageSchemaGenerateSftpPassword).(bool) {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.GenerateSftpPassword = true })
	}
	location := strings.TrimSpace(d.Get(StorageSchemaLocation).(string))
	if location != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Location = location })
	}
	name := strings.TrimSpace(d.Get(StorageSchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Name = name })
	}
	sftpPassword := strings.TrimSpace(d.Get(StorageSchemaSftpPassword).(string))
	if sftpPassword != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.SftpPassword = sftpPassword })
	}
	sType := strings.TrimSpace(d.Get(StorageSchemaType).(string))
	if sType != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Type = sType })
	}

	result, err := client.CreateStorage(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create storage: %v", err))
	}
	d.SetId(fmt.Sprintf("%d", result.ID))
	defer func() {
		dErr = resourceStorageRead(ctx, d, m)
	}()
	*id = int(result.ID)
	if result.Credentials.SftpPassword != "" {
		_ = d.Set(StorageSchemaSftpPassword, result.Credentials.SftpPassword)
	}
	if result.Credentials.S3.AccessKey != "" {
		_ = d.Set(StorageSchemaGenerateS3AccessKey, result.Credentials.S3.AccessKey)
	}
	if result.Credentials.S3.SecretKey != "" {
		_ = d.Set(StorageSchemaGenerateS3SecretKey, result.Credentials.S3.SecretKey)
	}

	keyId := d.Get(StorageSchemaKeyId).(int)
	if keyId == 0 {
		return dErr
	}
	keyOpts := []func(opt *storage.KeyLinkHTTPParams){
		func(opt *storage.KeyLinkHTTPParams) { opt.Context = ctx },
		func(opt *storage.KeyLinkHTTPParams) { opt.ID = result.ID; opt.KeyID = int64(keyId) },
	}
	err = client.LinkKeyToStorage(keyOpts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("link key to storage: %w", err))
	}

	return dErr
}

func resourceStorageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start Storage Resource reading (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Resource reading")
	if resourceId == "" {
		return diag.Errorf("get storage: empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	opts := []func(opt *storage.StorageListHTTPV2Params){
		func(opt *storage.StorageListHTTPV2Params) { opt.Context = ctx },
		func(opt *storage.StorageListHTTPV2Params) { opt.ID = &resourceId },
	}
	result, err := client.StoragesList(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("storages list: %w", err))
	}
	if len(result) != 1 {
		return diag.Errorf("get storage: wrong length of search result (%d), want 1", len(result))
	}
	st := result[0]

	nameParts := strings.Split(st.Name, "-")
	if len(nameParts) > 1 {
		clientID, _ := strconv.ParseInt(nameParts[0], 10, 64)
		_ = d.Set(StorageSchemaClientId, int(clientID))
		_ = d.Set(StorageSchemaName, strings.Join(nameParts[1:], "-"))
	} else {
		_ = d.Set(StorageSchemaName, st.Name)
	}
	_ = d.Set(StorageSchemaServerAlias, st.ServerAlias)
	_ = d.Set(StorageSchemaExpires, st.Expires)
	_ = d.Set(StorageSchemaId, st.ID)
	_ = d.Set(StorageSchemaType, st.Type)
	_ = d.Set(StorageSchemaLocation, st.Location)

	return nil
}

func resourceStorageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start Storage Resource updating (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Resource updating")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := make([]func(opt *storage.StorageUpdateHTTPParams), 0)
	opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Context = ctx })
	opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.ID = id })
	expires := strings.TrimSpace(d.Get(StorageSchemaExpires).(string))
	if expires != "" {
		opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Body.Expires = expires })
	}
	alias := strings.TrimSpace(d.Get(StorageSchemaServerAlias).(string))
	if alias != "" {
		opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Body.ServerAlias = alias })
	}
	_, err = client.ModifyStorage(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("update storage: %w", err))
	}

	return resourceStorageRead(ctx, d, m)
}

func resourceStorageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start Storage Resource deleting (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Resource deleting")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := []func(opt *storage.StorageDeleteHTTPParams){
		func(opt *storage.StorageDeleteHTTPParams) { opt.Context = ctx },
		func(opt *storage.StorageDeleteHTTPParams) { opt.ID = id },
	}
	err = client.DeleteStorage(opts...)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func storageResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		resourceID = fmt.Sprint(d.Get(StorageSchemaId).(int))
	}
	return resourceID
}

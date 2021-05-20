package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/G-Core/gcorelabs-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageSchemaGenerateSftpPassword = "generate_sftp_password"
	StorageSchemaLocation             = "location"
	StorageSchemaName                 = "name"
	StorageSchemaType                 = "type"
	StorageSchemaId                   = "id"
	StorageSchemaSftpPassword         = "sftp_password"
	StorageSchemaKeyId                = "link_key_id"
	StorageSchemaExpires              = "expires"
	StorageSchemaServerAlias          = "server_alias"
)

func resourceStorageResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A id of new storage resource.",
			},
			StorageSchemaName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name of new storage resource.",
			},
			StorageSchemaType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A type of new storage resource.",
			},
			StorageSchemaLocation: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A location of new storage resource.",
			},
			StorageSchemaSftpPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A sftp password for new storage resource.",
			},
			StorageSchemaGenerateSftpPassword: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "An auto generated sftp password for new storage resource.",
			},
			StorageSchemaKeyId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An key id to link with new storage resource.",
			},
		},
		CreateContext: resourceStorageResourceCreate,
		ReadContext:   resourceStorageResourceRead,
		UpdateContext: resourceStorageResourceUpdate,
		DeleteContext: resourceStorageResourceDelete,
		Description:   "Represent storage resource.",
	}
}

func resourceStorageResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (dErr diag.Diagnostics) {
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
		dErr = resourceStorageResourceRead(ctx, d, m)
	}()
	*id = int(result.ID)

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

func resourceStorageResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	sftpPass := ""
	if st.Credentials != nil {
		sftpPass = st.Credentials.SftpPassword
	}
	_ = d.Set(StorageSchemaServerAlias, st.ServerAlias)
	_ = d.Set(StorageSchemaExpires, st.Expires)
	_ = d.Set(StorageSchemaId, st.ID)
	_ = d.Set(StorageSchemaName, st.Name)
	_ = d.Set(StorageSchemaType, st.Type)
	_ = d.Set(StorageSchemaLocation, st.Location)
	_ = d.Set(StorageSchemaSftpPassword, sftpPass)

	return nil
}

func resourceStorageResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return resourceCDNResourceRead(ctx, d, m)
}

func resourceStorageResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		resourceID = strings.TrimSpace(d.Get(StorageSchemaId).(string))
	}
	return resourceID
}

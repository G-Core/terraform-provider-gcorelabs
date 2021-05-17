package gcore

import (
	"context"
	"fmt"
	"github.com/G-Core/gcorelabs-storage-sdk-go/swagger/client/key"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageKeySchemaKey  = "key"
	StorageKeySchemaName = "name"
	StorageKeySchemaId   = "id"
)

func resourceStorageKeyResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageKeySchemaName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name of new storage key resource.",
			},
			StorageKeySchemaKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A body of of new storage key resource.",
			},
		},
		CreateContext: resourceStorageKeyResourceCreate,
		ReadContext:   resourceStorageKeyResourceRead,
		DeleteContext: resourceStorageKeyResourceDelete,
		Description:   "Represent storage key resource.",
	}
}

func resourceStorageKeyResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (dErr diag.Diagnostics) {
	id := new(int)
	log.Println("[DEBUG] Start Storage Key Resource creating")
	defer log.Printf("[DEBUG] Finish Storage Key Resource creating (id=%d)\n", *id)
	config := m.(*Config)
	client := config.StorageClient

	opts := make([]func(opt *key.KeyCreateHTTPParams), 0)
	opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Context = ctx })
	name := strings.TrimSpace(d.Get(StorageKeySchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Body.Name = name })
	}
	keyText := strings.TrimSpace(d.Get(StorageKeySchemaKey).(string))
	if keyText != "" {
		opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Body.Key = keyText })
	}

	result, err := client.CreateKey(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create storage key: %v", err))
	}
	d.SetId(fmt.Sprintf("%d", result.ID))
	return resourceStorageKeyResourceRead(ctx, d, m)

}

func resourceStorageKeyResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageKeyResourceID(d)
	log.Printf("[DEBUG] Start Storage Key Resource reading (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Key Resource reading")
	if resourceId == "" {
		return diag.Errorf("get storage: empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	opts := []func(opt *key.KeyListHTTPV2Params){
		func(opt *key.KeyListHTTPV2Params) { opt.Context = ctx },
		func(opt *key.KeyListHTTPV2Params) { *opt.ID = resourceId },
	}
	name := strings.TrimSpace(d.Get(StorageKeySchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *key.KeyListHTTPV2Params) { *opt.Name = name })
	}
	result, err := client.KeysList(opts...)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) != 1 {
		return diag.Errorf("get storage key: wrong length of search result (%d), want 1", len(result))
	}
	st := result[0]

	_ = d.Set(StorageKeySchemaId, st.ID)
	_ = d.Set(StorageKeySchemaName, st.Name)

	return nil
}

func resourceStorageKeyResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageKeyResourceID(d)
	log.Printf("[DEBUG] Start Storage Key Resource deleting (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Key Resource deleting")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := []func(opt *key.KeyDeleteHTTPParams){
		func(opt *key.KeyDeleteHTTPParams) { opt.Context = ctx },
		func(opt *key.KeyDeleteHTTPParams) { opt.ID = id },
	}
	err = client.DeleteKey(opts...)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func storageKeyResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		resourceID = strings.TrimSpace(d.Get(StorageKeySchemaId).(string))
	}
	return resourceID
}

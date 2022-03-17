package gcore

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStorageS3Bucket() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageS3BucketSchemaStorageID: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "An id of existing storage resource.",
			},
			StorageS3BucketSchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					storageName := i.(string)
					if !regexp.MustCompile(`^[\w\-]+$`).MatchString(storageName) ||
						len(storageName) > 63 ||
						len(storageName) < 3 {
						return diag.Errorf("bucket name can't be empty and can have only letters & numbers. it also should be less than 63 symbols")
					}
					return nil
				},
				Description: "A name of storage bucket resource.",
			},
		},
		ReadContext: resourceStorageS3BucketRead,
		Description: "Represent storage s3 bucket resource.",
	}
}

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_region Data Source - terraform-provider-gcorelabs"
subcategory: ""
description: |-
  Represent region data
---

# gcore_region (Data Source)

Represent region data

## Example Usage

```terraform
provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Displayed region name

### Optional

- **id** (String) The ID of this resource.



provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage resource only and provide permanent_api_token without user_name & password
  gcore_platform = "https://api.gcdn.co"
  gcore_storage_api = "https://api.gcorelabs.com/storage"
}

resource "gcore_storage_s3_bucket" "example_s3_bucket" {
  name = "example1bucket2name"
  storage_id = 1
}

provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage resource only and provide permanent_api_token without user_name & password


}

data "gcore_storage_s3" "example_s3" {
  name = "example"
}

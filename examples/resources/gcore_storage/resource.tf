provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage resource only and provide permanent_api_token without user_name & password
  gcore_platform = "https://api.gcdn.co"
  gcore_storage_api = "https://storage.gcorelabs.com/api"
}

resource "gcore_storage" "tf_example_sftp" {
  name = "tf_example"
  location = "mia"
  type = "sftp"
  ssh_key_id = 199 // can be used for sftp type only
}

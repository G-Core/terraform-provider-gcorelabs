provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_storage_sftp_key" "example_key" {
  name = "example"
}

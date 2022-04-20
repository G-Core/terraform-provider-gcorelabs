provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_storage_s3" "example_s3" {
  name = "example"
}

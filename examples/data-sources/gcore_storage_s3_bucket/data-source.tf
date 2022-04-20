provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_storage_s3_bucket" "example_s3_bucket" {
  storage_id = 1
  name       = "example1bucket2name"
}

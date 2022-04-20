provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_storage_s3_bucket" "example_s3_bucket" {
  name       = "example1bucket2name"
  storage_id = 1
}

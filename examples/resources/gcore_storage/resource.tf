provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_storage_api = "https://storage.gcorelabs.com/api"
}

resource "gcore_storage" "tf_example_s3" {
  name = "tf_example"
  location = "s-ed1"
  type = "s3"
  link_key_id = 199
}

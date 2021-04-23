provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_servergroup" "default" {
  name = "default"
  policy = "affinity"
  region_id = 1
  project_id = 1
}

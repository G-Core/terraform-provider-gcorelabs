provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_snapshot" "snapshot" {
  project_id = 1
  region_id = 1
  name = "snapshot example"
  volume_id = "28e9edcb-1593-41fe-971b-da729c6ec301"
  description = "snapshot example description"
}



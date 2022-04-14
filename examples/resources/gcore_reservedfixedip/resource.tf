provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = 1
  region_id = 1
  type = "external"
  is_vip = false
}
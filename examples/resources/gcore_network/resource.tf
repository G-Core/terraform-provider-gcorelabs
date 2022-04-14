provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_network" "network" {
  name = "network_example"
  mtu = 1450
  type = "vxlan"
  region_id = 1
  project_id = 1
}
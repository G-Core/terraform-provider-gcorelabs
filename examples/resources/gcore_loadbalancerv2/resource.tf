provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_loadbalancerv2" "lb" {
  project_id = 1
  region_id = 1
  name = "test"
  flavor = "lb1-1-2"
}
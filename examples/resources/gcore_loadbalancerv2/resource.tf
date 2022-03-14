provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_loadbalancerv2" "lb" {
  project_id = 1
  region_id = 1
  name = "test"
  flavor = "lb1-1-2"
}
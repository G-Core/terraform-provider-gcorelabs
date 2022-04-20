provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_loadbalancerv2" "lb" {
  project_id = 1
  region_id  = 1
  name       = "test"
  flavor     = "lb1-1-2"
}
provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_network" "network" {
  name       = "network_example"
  type       = "vxlan"
  region_id  = 1
  project_id = 1
}

provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_network" "network" {
  name = "network_example"
  mtu = 1450
  type = "vxlan"
  region_id = 1
  project_id = 1
}
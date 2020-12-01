provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_volume" "volume" {
  name = "volume_example"
  type_name = "standard"
  size = 1
  region_id = 1
  project_id = 1
}
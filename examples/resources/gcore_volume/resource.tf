provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_volume" "volume" {
  name = "volume_example"
  type_name = "standard"
  size = 1
  region_id = 1
  project_id = 1
}
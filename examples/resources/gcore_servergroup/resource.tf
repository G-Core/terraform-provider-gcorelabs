provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_servergroup" "default" {
  name = "default"
  policy = "affinity"
  region_id = 1
  project_id = 1
}

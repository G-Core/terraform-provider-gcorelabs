provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_snapshot" "snapshot" {
  project_id = 1
  region_id = 1
  name = "snapshot example"
  volume_id = "28e9edcb-1593-41fe-971b-da729c6ec301"
  description = "snapshot example description"
  metadata = {
    env = "test"
  }
}



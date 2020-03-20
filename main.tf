provider "gcore" {
  username = "vvayner@drova.io"
  password = "P@ssw0rd1!"
}

resource "gcore_volume" "simple_volume" {
  name = 156
  size = 2
  type_name = "ssd_hiiops"
  region_id = 1
  project_id = 81
  source = "new-volume"
  #count = 2
}

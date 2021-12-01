provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_loadbalancer" "lb" {
  project_id = 1
  region_id = 1
  name = "test"
  flavor = "lb1-1-2"
  //when upgrading to version 0.2.28 nested listener max length reduced to 1
  //that mean, if you had more than one nested listener and removed them from
  //schema they not delete in the cloud. User has to delete it manually and
  //recreate as gcore_lblistener resource
  listener {
    name = "test"
    protocol = "HTTP"
    protocol_port = 80
  }
}
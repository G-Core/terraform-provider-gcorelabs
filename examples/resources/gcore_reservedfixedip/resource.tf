provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = 1
  region_id = 1
  type = "external"
  is_vip = false
}
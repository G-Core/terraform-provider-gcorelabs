provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_ddos_protection" "ddos_protection" {
  project_id = 1
  region_id = 1
  profile_template = 63
  ip_address = "10.94.77.72"
  bm_instance_id = "99cd3a2d-607f-4fbb-91d9-01fe926b1e7f"
  fields {
    base_field = 118
    field_value = [33033]
  }
}
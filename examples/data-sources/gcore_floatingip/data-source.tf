provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

data "gcore_project" "pr" {
  name = "test"
}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}

data "gcore_floatingip" "ip" {
  floating_ip_address = "10.100.179.172"
  region_id = data.gcore_region.rg.id
  project_id = data.gcore_project.pr.id
}

output "view" {
  value = data.gcore_floatingip.ip
}


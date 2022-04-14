provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_network" "network" {
  name = "network_example"
  mtu = 1450
  type = "vxlan"
  region_id = 1
  project_id = 1
}

resource "gcore_subnet" "subnet" {
  name = "subnet_example"
  cidr = "192.168.10.0/24"
  network_id = gcore_network.network.id
  dns_nameservers = var.dns_nameservers

  dynamic host_routes {
    iterator = hr
    for_each = var.host_routes
      content {
        destination = hr.value.destination
        nexthop = hr.value.nexthop
      }
  }

  gateway_ip = "192.168.10.1"
  region_id = 1
  project_id = 1
}
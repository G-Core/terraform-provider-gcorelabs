provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
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
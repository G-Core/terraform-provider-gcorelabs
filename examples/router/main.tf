provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_router" "router" {
  name = "router_example"

  dynamic external_gateway_info {
  iterator = egi
  for_each = var.external_gateway_info
  content {
    type = egi.value.type
    enable_snat = egi.value.enable_snat
    network_id = egi.value.network_id
    }
  }

  dynamic interfaces {
  iterator = ifaces
  for_each = var.interfaces
  content {
    type = ifaces.value.type
    subnet_id = ifaces.value.subnet_id
    }
  }

  dynamic routes {
  iterator = rs
  for_each = var.routes
  content {
    destination = rs.value.destination
    nexthop = rs.value.nexthop
    }
  }

  region_id = 1
  project_id = 1
}
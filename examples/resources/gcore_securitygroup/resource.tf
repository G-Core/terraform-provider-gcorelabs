provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_securitygroup" "sg" {
  name = "test sg"
  region_id = 1
  project_id = 1

  security_group_rules {
    direction = "egress"
    ethertype = "IPv4"
    protocol = "tcp"
    port_range_min = 19990
    port_range_max = 19990
  }

  security_group_rules {
    direction = "ingress"
    ethertype = "IPv4"
    protocol = "tcp"
    port_range_min = 19990
    port_range_max = 19990
  }

  security_group_rules {
    direction = "egress"
    ethertype = "IPv4"
    protocol = "vrrp"
  }
}

provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
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
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]

  host_routes {
    destination = "10.0.3.0/24"
    nexthop = "10.0.0.13"
  }

  gateway_ip = "192.168.10.1"
  region_id = 1
  project_id = 1
}

resource "gcore_volume" "first_volume" {
  name = "boot volume"
  type_name = "ssd_hiiops"
  size = 5
  image_id = "f4ce3d30-e29c-4cfd-811f-46f383b6081f"
  region_id = 1
  project_id = 1
}

resource "gcore_volume" "second_volume" {
  name = "second volume"
  type_name = "ssd_hiiops"
  size = 5
  region_id = 1
  project_id = 1
}

resource "gcore_instance" "instance" {
  flavor_id = "g1-standard-2-4"
  name = "test"

  volume {
    source = "existing-volume"
    volume_id = gcore_volume.first_volume.id
    boot_index = 0
  }

  volume {
    source = "existing-volume"
    volume_id = gcore_volume.second_volume.id
    boot_index = 1
  }

  interface {
    type = "subnet"
    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
    //port_id = null
    //ip_address = null
    //fip_source = null
    //existing_fip_id = null
  }


  security_group {
    id = "d75db0b2-58f1-4a11-88c6-a932bb897310"
    name = "default"
  }

  metadata {
    key = "some_key"
    value = "some_data"
  }

  configuration {
    key = "some_key"
    value = "some_data"
  }

  region_id = 1
  project_id = 1
}



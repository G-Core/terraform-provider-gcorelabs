terraform {
  required_version = ">= 0.13.0"
  required_providers {
    gcore = {
      source  = "G-Core/gcorelabs"
      version = "0.1.9"
    }
  }
}

provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_keypair" "kp" {
  project_id = 1
  public_key = "your oub key"
  sshkey_name = "testkey"
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

resource "gcore_subnet" "subnet2" {
  name = "subnet2_example"
  cidr = "192.168.20.0/24"
  network_id = gcore_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]

  host_routes {
    destination = "10.0.3.0/24"
    nexthop = "10.0.0.13"
  }

  gateway_ip = "192.168.20.1"
  region_id = 1
  project_id = 1
}

resource "gcore_volume" "first_volume" {
  name = "boot volume"
  type_name = "ssd_hiiops"
  size = 6
  image_id = "f4ce3d30-e29c-4cfd-811f-46f383b6081f"
  region_id = 1
  project_id = 1
}

resource "gcore_volume" "second_volume" {
  name = "second volume"
  type_name = "ssd_hiiops"
  image_id = "f4ce3d30-e29c-4cfd-811f-46f383b6081f"
  size = 6
  region_id = 1
  project_id = 1
}

resource "gcore_volume" "third_volume" {
  name = "third volume"
  type_name = "ssd_hiiops"
  size = 6
  region_id = 1
  project_id = 1
}

resource "gcore_instance" "instance" {
  flavor_id = "g1-standard-2-4"
  name = "test"
  keypair_name = gcore_keypair.kp.sshkey_name

  volume {
    source = "existing-volume"
    volume_id = gcore_volume.first_volume.id
    boot_index = 0
  }

  interface {
    type = "subnet"
    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }

  interface {
    type = "subnet"
     network_id = gcore_network.network.id
     subnet_id = gcore_subnet.subnet2.id
  }

  security_group {
    id = "66988147-f1b9-43b2-aaef-dee6d009b5b7"
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

resource "gcore_loadbalancer" "lb" {
  project_id = 1
  region_id = 1
  name = "test1"
  flavor = "lb1-1-2"
  listener {
    name = "test"
    protocol = "HTTP"
    protocol_port = 80
  }
}

resource "gcore_lbpool" "pl" {
  project_id = 1
  region_id = 1
  name = "test_pool1"
  protocol = "HTTP"
  lb_algorithm = "LEAST_CONNECTIONS"
  loadbalancer_id = gcore_loadbalancer.lb.id
  listener_id = gcore_loadbalancer.lb.listener.0.id
  health_monitor {
    type = "PING"
    delay = 60
    max_retries = 5
    timeout = 10
  }
  session_persistence {
    type = "APP_COOKIE"
    cookie_name = "test_new_cookie"
  }
}

resource "gcore_lbmember" "lbm" {
  project_id = 1
  region_id = 1
  pool_id = gcore_lbpool.pl.id
  instance_id = gcore_instance.instance.id
  address = tolist(gcore_instance.instance.interface).0.ip_address
  protocol_port = 8081
  weight = 5
}

resource "gcore_instance" "instance2" {
  flavor_id = "g1-standard-2-4"
  name = "test2"
  keypair_name = gcore_keypair.kp.sshkey_name

  volume {
    source = "existing-volume"
    volume_id = gcore_volume.second_volume.id
    boot_index = 0
  }

  volume {
  	source = "existing-volume"
  	volume_id = gcore_volume.third_volume.id
  	boot_index = 1
  }

  interface {
    type = "subnet"
    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }

  security_group {
    id = "66988147-f1b9-43b2-aaef-dee6d009b5b7"
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

resource "gcore_lbmember" "lbm2" {
  project_id = 1
  region_id = 1
  pool_id = gcore_lbpool.pl.id
  instance_id = gcore_instance.instance2.id
  address = tolist(gcore_instance.instance2.interface).0.ip_address
  protocol_port = 8081
  weight = 5
}
provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_volume" "first_volume" {
  name = "boot volume"
  type_name = "ssd_hiiops"
  size = 5
  image_id = "f3847215-e4d7-4e64-8e69-14637e68e27f"
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

locals {
  volumes_ids = [gcore_volume.first_volume.id, gcore_volume.second_volume.id]
}

resource "gcore_instance" "instance" {
  flavor_id = "g1-standard-2-4"
  name = var.names

  dynamic volumes {
  iterator = vol
  for_each = local.volumes_ids
  content {
    boot_index = index(local.volumes_ids, vol.value)
    source = "existing-volume"
    volume_id = vol.value
    }
  }

  dynamic interfaces {
  iterator = iface
  for_each = var.interfaces
  content {
    type = iface.value.type
    network_id = iface.value.network_id
    subnet_id = iface.value.subnet_id
    fip_source = iface.value.fip_source
    existing_fip_id =iface.value.existing_fip_id
    port_id = iface.value.port_id
    ip_address = iface.value.ip_address
    }
  }

  dynamic security_groups {
  iterator = sg
  for_each = var.security_groups
  content {
    id = sg.value.id
    name = sg.value.name
    }
  }

  dynamic metadata {
  iterator = md
  for_each = var.metadata
  content {
    key = md.value.key
    value = md.value.value
    }
  }

  dynamic configuration {
  iterator = cfg
  for_each = var.configuration
  content {
    key = cfg.value.key
    value = cfg.value.value
    }
  }

  region_id = 1
  project_id = 1
}
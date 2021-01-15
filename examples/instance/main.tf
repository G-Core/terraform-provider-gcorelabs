provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_instance" "instance" {
  flavor_id = "g1-standard-2-4"
  name = var.names

  dynamic volumes {
  iterator = vol
  for_each = var.volumes
  content {
    source = vol.value.source
    type_name = vol.value.type_name
    size = vol.value.size
    name = vol.value.name
    boot_index = vol.value.boot_index
    image_id = vol.value.image_id
    }
  }

  dynamic interfaces {
  iterator = iface
  for_each = var.interfaces
  content {
    type = iface.value.type
    network_id = iface.value.network_id
    subnet_id = iface.value.subnet_id
    }
  }

  dynamic security_groups {
  iterator = sg
  for_each = var.security_groups
  content {
    id = sg.value.id
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
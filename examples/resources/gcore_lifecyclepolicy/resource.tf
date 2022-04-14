provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_lifecyclepolicy" "lp" {
  project_id = 1
  region_id = 1
  name = "test"
  status = "active"
  action = "volume_snapshot"
  volume {
    id = "fe93bfdd-4ce3-4041-b89b-4f10d0d49498"
  }
  schedule {
    max_quantity = 4
    interval {
      weeks = 1
      days = 2
      hours = 3
      minutes = 4
    }
    resource_name_template = "reserve snap of the volume {volume_id}"
    retention_time {
      weeks = 4
      days = 3
      hours = 2
      minutes = 1
    }
  }
}
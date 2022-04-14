provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cloud_api = "https://api.gcorelabs.com/cloud"
}

resource "gcore_k8s" "v" {
  project_id = 1
  region_id = 1
  name = "tf-k8s"
  fixed_network = "6bf878c1-1ce4-47c3-a39b-6b5f1d79bf25"
  fixed_subnet = "dc3a3ea9-86ae-47ad-a8e8-79df0ce04839"
  pool {
    name = "tf-pool"
    flavor_id = "g1-standard-1-2"
    min_node_count = 1
    max_node_count = 2
    node_count = 1
    docker_volume_size = 2
  }
}


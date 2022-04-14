provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_k8s_pool" "v" {
  project_id = 1
  region_id = 1
  cluster_id = "6bf878c1-1ce4-47c3-a39b-6b5f1d79bf25"
  name = "tf-pool"
  flavor_id = "g1-standard-1-2"
  min_node_count = 1
  max_node_count = 2
  node_count = 1
  docker_volume_size = 2
}


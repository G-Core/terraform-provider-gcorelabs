provider gcore {
  user_name = "test"
  password = "test"

}

data "gcore_k8s" "v" {
  project_id = 1
  region_id = 1
  cluster_id = "dc3a3ea9-86ae-47ad-a8e8-79df0ce04839"
}


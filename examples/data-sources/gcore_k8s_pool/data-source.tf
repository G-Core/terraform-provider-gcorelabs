provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

data "gcore_k8s_pool" "v" {
  project_id = 1
  region_id = 1
  cluster_id = "6bf878c1-1ce4-47c3-a39b-6b5f1d79bf25"
  pool_id = "dc3a3ea9-86ae-47ad-a8e8-79df0ce04839"
}


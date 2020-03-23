resource "gcore_volumeV1" "volume_example" {
  name = "volume_example"
  size = 1
  type_name = "standard"
  source = "new-volume"
  region_id = 1
  project_id = 1
}
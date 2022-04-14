provider gcore {
  user_name = "test"
  password = "test"


}

data "gcore_project" "pr" {
  name = "test"
}

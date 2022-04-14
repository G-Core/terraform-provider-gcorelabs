provider gcore {
  user_name = "test"
  password = "test"


}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}

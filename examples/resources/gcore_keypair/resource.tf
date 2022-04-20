provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_keypair" "kp" {
  project_id  = 1
  public_key  = "your public key here"
  sshkey_name = "test"
}

output "kp" {
  value = gcore_keypair.kp
}

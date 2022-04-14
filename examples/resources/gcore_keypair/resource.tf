provider gcore {
  user_name = "test"
  password = "test"


}

resource "gcore_keypair" "kp" {
  project_id = 1
  public_key = "your public key here"
  sshkey_name = "test"
}

output "kp" {
  value = gcore_keypair.kp
}

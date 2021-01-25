provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_keypair" "kp" {
  project_id = 1
  public_key = "your public key here"
  sshkey_name = "test"
}

output "kp" {
  value = gcore_keypair.kp
}

provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}

variable "cert" {
  type = string
  sensitive = true
}

variable "private_key" {
  type = string
  sensitive = true
}

resource "gcore_cdn_sslcert" "cdnopt_cert" {
  name = "Test cert for cdnopt_bookatest_by"
  cert = var.cert
  private_key = var.private_key
}


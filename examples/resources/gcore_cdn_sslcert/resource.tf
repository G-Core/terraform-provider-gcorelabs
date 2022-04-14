provider gcore {
  # G-Core dashboard => Profile => API tokens => Create token
  permanent_api_token = ""

  # user_name = "test"
  # password = "test"
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


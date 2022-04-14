provider gcore {
  # G-Core dashboard => Profile => API tokens => Create token
  permanent_api_token = ""

  # user_name = "test"
  # password = "test"

  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_cdn_api = "https://api.gcorelabs.com"
}

resource "gcore_cdn_origingroup" "origin_group_1" {
  name = "origin_group_1"
  use_next = true
  origin {
    source = "example.com"
    enabled = true
  }
  origin {
    source = "mirror.example.com"
    enabled = true
    backup = true
  }
}

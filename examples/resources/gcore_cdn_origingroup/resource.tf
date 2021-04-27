provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}

resource "gcore_cdn_origingroup" "origin_group_1" {
  name = "origin_group_1"
  use_next = true
  origin {
    source = "example.com"
    enabled = false
  }
  origin {
    source = "mirror.example.com"
    enabled = true
    backup = true
  }
}

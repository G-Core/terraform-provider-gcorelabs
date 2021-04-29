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

resource "gcore_cdn_resource" "cdn_example_com" {
  cname = "cdn.example.com"
  origin_group = gcore_cdn_origingroup.origin_group_1.id
  origin_protocol = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]
}

resource "gcore_cdn_rule" "cdn_example_com_rule_1" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name = "All images"
  rule = "/folder/images/*.png"
  rule_type = 0
}

resource "gcore_cdn_rule" "cdn_example_com_rule_2" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name = "All scripts"
  rule = "/folder/images/*.js"
  rule_type = 0
}

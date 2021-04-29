provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}


resource "gcore_cdn_resource" "cdn_example_com" {
  cname = "cdn.example.com"
  origin_group = 11
  origin_protocol = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]
}

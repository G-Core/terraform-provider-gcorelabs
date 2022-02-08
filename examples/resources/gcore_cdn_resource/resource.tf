provider gcore {
  # G-Core dashboard => Profile => API tokens => Create token
  permanent_api_token = ""

  # user_name = "test"
  # password = "test"

  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}


resource "gcore_cdn_resource" "cdn_example_com" {
  cname = "cdn.example.com"
  origin_group = 11
  origin_protocol = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]

  options {
    browser_cache_settings {
      value = "1d"
    }
    redirect_http_to_https {
      value = true
    }
    gzip_on {
      value = true
    }
    cors {
      value = [
        "*"
      ]
    }
    rewrite {
      body = "/(.*) /$1"
    }
    webp {
      jpg_quality = 55
      png_quality = 66
    }
  }
}

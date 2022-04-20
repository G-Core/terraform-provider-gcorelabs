provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}


resource "gcore_cdn_resource" "cdn_example_com" {
  cname               = "cdn.example.com"
  origin_group        = gcore_cdn_origingroup.origin_group_1.id
  origin_protocol     = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]

  options {
    edge_cache_settings {
      default = "8d"
    }
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

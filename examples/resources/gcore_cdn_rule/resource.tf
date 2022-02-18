provider gcore {
  # G-Core dashboard => Profile => API tokens => Create token
  permanent_api_token = ""

  # user_name = "test"
  # password = "test"

  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}

resource "gcore_cdn_rule" "cdn_example_com_rule_1" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name = "All PNG images"
  rule = "/folder/images/*.png"
  rule_type = 0

  options {
    edge_cache_settings {
      default = "14d"
    }
    browser_cache_settings {
      value = "14d"
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
    ignore_query_string {
      value = true
    }
  }
}

resource "gcore_cdn_rule" "cdn_example_com_rule_2" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name = "All JS scripts"
  rule = "/folder/images/*.js"
  rule_type = 0
  origin_protocol = "HTTP"

  options {
    redirect_http_to_https {
      enabled = false
      value = true
    }
    gzip_on {
      enabled = false
      value = true
    }
    query_params_whitelist {
      value = [
        "abc",
      ]
    }
  }
}

resource "gcore_cdn_origingroup" "origin_group_1" {
  name = "origin_group_1"
  use_next = true
  origin {
    source = "example.com"
    enabled = true
  }
}

resource "gcore_cdn_resource" "cdn_example_com" {
  cname = "cdn.example.com"
  origin_group = gcore_cdn_origingroup.origin_group_1.id
  origin_protocol = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]
}

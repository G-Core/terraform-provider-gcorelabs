provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_loadbalancer" "lb" {
  project_id = 1
  region_id  = 1
  name       = "test1"
  flavor     = "lb1-1-2"
  listener {
    name          = "test"
    protocol      = "HTTP"
    protocol_port = 80
  }
}

resource "gcore_lbpool" "pl" {
  project_id      = 1
  region_id       = 1
  name            = "test_pool1"
  protocol        = "HTTP"
  lb_algorithm    = "LEAST_CONNECTIONS"
  loadbalancer_id = gcore_loadbalancer.lb.id
  listener_id     = gcore_loadbalancer.lb.listener.0.id
  health_monitor {
    type        = "PING"
    delay       = 60
    max_retries = 5
    timeout     = 10
  }
  session_persistence {
    type        = "APP_COOKIE"
    cookie_name = "test_new_cookie"
  }
}

provider gcore {
  user_name = "test@test.test"
  password = "testtest"
  gcore_platform = "http://api.stg-45.staging.gcdn.co"
  gcore_api = "http://10.100.179.92:33081"
}

resource "gcore_loadbalancer" "lb" {
  project_id = 1
  region_id = 1
  name = "test1"
  flavor = "lb1-1-2"
  listeners {
    name = "test"
    protocol = "HTTP"
    protocol_port = 80
  }
}

resource "gcore_lbpool" "pl" {
  project_id = 1
  region_id = 1
  name = "test_pool1"
  protocol = "HTTP"
  lb_algorithm = "LEAST_CONNECTIONS"
  loadbalancer_id = gcore_loadbalancer.lb.id
  listener_id = gcore_loadbalancer.lb.listeners.0.id
  health_monitor {
    type = "PING"
    delay = 60
    max_retries = 5
    timeout = 10
  }
  session_persistence {
    type = "APP_COOKIE"
    cookie_name = "test_new_cookie"
  }
}

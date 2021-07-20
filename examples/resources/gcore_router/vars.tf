variable "external_gateway_info" {
  type = list(object({
    type = string
    enable_snat = bool
    network_id = string
  }))
  default = [
    {
      type = "manual"
      enable_snat = false
      network_id = "" //set external network id
    },
  ]
}

variable "interfaces" {
  type = list(object({
    type = string
    subnet_id = string
  }))
  default = [
    {
      type = "subnet"
      subnet_id = "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
    },
    {
      type = "subnet"
      subnet_id = "f3f6a294-a319-4db4-84b6-6016a3481924"
    },
  ]
}

variable "routes" {
  type = list(object({
    destination = string
    nexthop = string
  }))
  default = [
    {
      destination = "192.168.101.0/24"
      nexthop = "192.168.100.2"
    },
    {
      destination = "192.168.102.0/24"
      nexthop = "192.168.100.3"
    },
  ]
}
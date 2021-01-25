variable "names" {
  type = list
  default = ["instance_example"]
}

/*
detach interface:
... = null
port_id = "cd69fa45-688c-4c87-9ee2-1684f27774ad"
ip_address = "192.168.55.118"
... = null
attach interface:
... = null
type = "subnet"
subnet_id = "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
... = null
*/
variable "interfaces" {
  type = list(object({
    type = string
    subnet_id = string
    network_id = string
    fip_source = string
    existing_fip_id = string
    port_id = string
    ip_address = string
  }))
  default = [
    {
      type = "subnet"
      network_id = "900f9c9f-35db-40a5-88ce-efb773947c0a"
      subnet_id = "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
      port_id = null
      ip_address = null
      fip_source = null
      existing_fip_id = null
    },
  ]
}

variable "security_groups" {
  type = list(object({
    id = string
    name = string
  }))
  default = [
    {
      id = "81e6dfd9-b646-4a5f-9064-cc224c63e545"
      name = "default"
    },
  ]
}

variable "metadata" {
  type = list(object({
    key = string
    value = string
  }))
  default = [
    {
      key = "some_key"
      value = "some_data"
    },
  ]
}

variable "configuration" {
  type = list(object({
    key = string
    value = string
  }))
  default = [
    {
      key = "some_key"
      value = "some_data"
    },
  ]
}
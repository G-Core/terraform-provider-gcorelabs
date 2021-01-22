variable "names" {
  type = list
  default = ["instance_example"]
}

variable "interfaces" {
  type = list(object({
    type = string
    subnet_id = string
    network_id = string
    port_id = string
    floating_ip = object({
                    source = string
                    existing_floating_id = string
                    })
  }))
  default = [
    {
      type = "subnet"
      network_id = "900f9c9f-35db-40a5-88ce-efb773947c0a"
      subnet_id = "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
      port_id = null
      floating_ip = null
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
      name = null
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
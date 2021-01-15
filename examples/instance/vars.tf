variable "names" {
  type = list
  default = ["instance_example"]
}

variable "volumes" {
  type = list(object({
    source = string
    boot_index = number
    type_name = string
    size = number
    name = string
    image_id = string
  }))
  default = [
    {
      source = "image"
      type_name = "ssd_hiiops"
      size = 5
      name = "boot volume"
      boot_index = 0
      image_id = "f3847215-e4d7-4e64-8e69-14637e68e27f"
    },
    {
      source = "new-volume"
      type_name = "ssd_hiiops"
      size = 5
      name = "empty volume"
      boot_index = 1
      image_id = null
    },
  ]
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
  }))
  default = [
    {
      id = "81e6dfd9-b646-4a5f-9064-cc224c63e545"
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
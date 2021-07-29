variable "stringList" {
  type = list(string)
  default = ["one", "two", "three"]
}

variable "intList" {
  type = list(number)
  default = [1, 2, 3]
}

variable "floatList" {
  type = list(number)
  default = [1.01, 2.02, 3.03]
}

variable "boolList" {
  type = list(bool)
  default = [true, true, false, true, false]
}

variable "list_no_type" {
  default = [1,2]
}

variable "set_var" {
  type = set(string)
  default = ["first","second"]
}

variable "tuple_var" {
  type = tuple([string, number, bool])
  default = ["one", 1, true]
}

variable "list_of_tuple" {
  type = list(tuple([string, number, bool]))
  default = [
    ["one", 1, true],
    ["two", 2, false]
    ]
}

variable "tuple_var_complex" {
  type = tuple([number, object({
    field1 = number
    field2 = number
  })])

  default = [10, {
    field1 = 11
    field2 = 12
  }]
}

variable "object_var" {
  type = object({
    name    = string
    address = string
  })
  default = {
      name = "pankaj"
      address = "pune"
  }
}

variable "map_var" {
  type = map
  default = {
    "5USD"  = "1xCPU-1GB"
    "10USD" = "1xCPU-2GB"
    "20USD" = "2xCPU-4GB"
  }
}

variable "map_var_complex" {
  type = map(object({
    name = string
    ID = number
  }))
  default = {
    "first"  = {
      name = "Thor"
      ID = 1
    }
    "second" = {
      name = "Antman"
      ID = 2
    }
  }
}

variable "object_list" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    },
    {
      internal = 4000
      external = 3000
      protocol = "udp"
    }
  ]
}

variable "object_list_complex" {
  type = list(object({
    key1 = list(number)
    key2 = object({
      port = number
    })
    key3 = map(string)
    key4 = map(number)
  }))
  default = [
    {
      key1 = [1, 2, 3]
      key2 = {
        port = 9010
      }
      key3 = {
        "name": "hero"
      }
      key4 = {
        "first": 11.23
        "second": 50
      }
    }
  ]
}
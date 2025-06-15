variable "project_env" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "table_name" { type = string }

variable "hash_key" {
  type    = string
  default = "pk"
}

variable "range_key" {
  type    = string
  default = null
}

variable "attributes" { type = list(object({ name = string, type = string })) }

variable "read_capacity" {
  type    = number
  default = 1
}

variable "write_capacity" {
  type    = number
  default = 1
}

variable "billing_mode" {
  type    = string
  default = "PAY_PER_REQUEST"
}

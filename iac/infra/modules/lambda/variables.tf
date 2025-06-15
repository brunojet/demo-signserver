variable "project_env" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "name" { type = string }
variable "handler" { type = string }
variable "runtime" { type = string }
variable "filename" { type = string }
variable "environment_variables" {
  type    = map(string)
  default = {}
}
variable "timeout" {
  type    = number
  default = 3
}
variable "memory_size" {
  type    = number
  default = 128
}

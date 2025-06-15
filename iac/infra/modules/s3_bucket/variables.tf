variable "project_env" {
  type = string
}

variable "tags" {
  type = map(string)
}


variable "bucket_name" { type = string }

variable "versioning" {
  type    = bool
  default = true
}

variable "project_env" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "name" { type = string }
variable "event_pattern" { type = string }
variable "target_arn" { type = string }

# Permite customizar o nome do bus, mas usa o default se não informado
variable "event_bus_name" {
  type        = string
  default     = "default"
  description = "Nome do EventBridge bus (default = 'default')"
}

terraform {
  required_providers {
    edgio = {
      source = "Edgio/edgio"
      version = "0.1.0"
    }
  }
}

variable "client_id" { type = string }
variable "client_secret" {  type = string }
variable "organization_id" { type = string }

provider "edgio" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

resource "edgio_property" "my_property" {
  organization_id = var.organization_id
  slug            = "edgio-environment-example"
}

resource "edgio_environment" "my_env" {
  property_id                 = edgio_property.my_property.id
  name                        = "main"
  only_maintainers_can_deploy = true  
  http_request_logging        = true
}

output "added_env" {
  value = edgio_environment.my_env
}
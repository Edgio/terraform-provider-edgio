terraform {
  required_providers {
    edgio = {
      source = "hashicorp.com/edu/edgio"
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
  slug = "edgio-property-example"
}

output "added_property" {
  value = edgio_property.my_property
}
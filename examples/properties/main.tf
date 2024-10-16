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

data "edgio_properties" "my_properties" {
  item_count = 100
  organization_id = var.organization_id
}

output "properties" {
  value = data.edgio_properties.my_properties.properties
}

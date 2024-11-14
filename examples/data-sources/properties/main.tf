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
variable "property_id" { type = string }

provider "edgio" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

data "edgio_properties" "my_properties" {
  item_count = 100
  property_id = var.property_id
}

output "properties" {
  value = data.edgio_properties.my_properties.properties
}

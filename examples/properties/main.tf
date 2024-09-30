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
   organization_id = var.organization_id
}

output "properties" {
  value = [for property in data.edgio_properties.my_properties.properties : {
    id: property.id,
    type: property.type,
    slug: property.slug,
    organization_id: property.organization_id,
    created_at: property.created_at,
    updated_at: property.updated_at,
  }]
}

output "item_count" {
  value = data.edgio_properties.my_properties.item_count
}

output "links" {
  value = {
    first: data.edgio_properties.my_properties.links.first.href,
    next: data.edgio_properties.my_properties.links.next.href,
    previous: data.edgio_properties.my_properties.links.previous.href,
    last: data.edgio_properties.my_properties.links.last.href,
  }
}

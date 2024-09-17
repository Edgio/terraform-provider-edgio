terraform {
  required_providers {
    edgio = {
      source = "hashicorp.com/edu/edgio"
    }
  }
}

provider "edgio" {
  client_id     = "f8c1d12a-ee43-44d9-816a-bd73b7441ca5"
  client_secret = "veBWKIS5vY9akbw5UaqksF7Et29lQnDo"
}

data "edgio_properties" "my_properties" {
   organization_id = "6b1e0c15-d302-4775-b731-efaa22b96617"
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

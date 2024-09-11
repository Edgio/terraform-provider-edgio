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

data "edgio_create_property" "my_property" {
  organization_id = "6b1e0c15-d302-4775-b731-efaa22b96617"
  slug = "my-new-property"
}

output "property_details" {
   value = data.edgio_create_property.my_property.slug
}
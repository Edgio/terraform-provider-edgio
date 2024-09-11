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

data "edgio_property" "my_specific_property" {
   property_id = "23e376ec-5aab-4a46-84d4-a15571b3c994"
}

output "property_details" {
  value = {
    id: data.edgio_property.my_specific_property.id,
    organization_id: data.edgio_property.my_specific_property.organization_id,
    slug: data.edgio_property.my_specific_property.slug,
    created_at: data.edgio_property.my_specific_property.created_at,
    updated_at: data.edgio_property.my_specific_property.updated_at,
    
  }
}
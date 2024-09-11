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

data "edgio_properties" "all_properties" {
   organization_id = "6b1e0c15-d302-4775-b731-efaa22b96617"
}

output "formatted_properties" {
  value = join("\n", [
    for p in data.edgio_properties.all_properties.properties : format(
      "Property Name: %s\n    ID: %s\n    Type: %s\n    Slug: %s\n    Organization ID: %s\n    Links:\n        First: %s\n        Next: %s\n        Previous: %s\n        Last: %s\n",
      p.name,
      p.id,
      p.type,
      p.slug,
      p.organization_id,
      p.links.first.href,
      p.links.next.href,
      p.links.previous.href,
      p.links.last.href
    )
  ])
}

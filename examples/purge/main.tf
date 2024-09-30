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

resource "edgio_purge_cache" "my_purge_cache" {
  environment_id = "6939fe34-a9aa-43e2-b2ff-9adc7f28a0df"
  purge_type     = "all_entries"
  values         =  []
}

output "added_purge_cache" {
  value = {
    percentage: edgio_purge_cache.my_purge_cache.progress_percentage
  }
}
 
 
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

resource "edgio_environment" "my_env" {
  property_id = "aba651c5-4bf5-426d-ad33-cf44b8aac63e"
  name        = "new-env-changed"
  can_members_deploy = true
  only_maintainers_can_deploy = false
  http_request_logging = true
}

output "added_env" {
  value = edgio_environment.my_env
}
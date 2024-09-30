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

data "edgio_tls_certs" "my_certs" {
   environment_id = "6939fe34-a9aa-43e2-b2ff-9adc7f28a0df"
   page = 1
   page_size = 10
}

output "all_my_certs" {
  value = data.edgio_tls_certs.my_certs
}

terraform {
  required_providers {
    edgio = {
      source = "hashicorp.com/edu/edgio"
    }
  }
}

variable "client_id" { type = string }
variable "client_secret" {  type = string }
variable "environment_id" { type = string }

provider "edgio" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

data "edgio_tls_certs" "my_certs" {
   environment_id = var.environment_id
   page = 1
   page_size = 10
}

output "all_my_certs" {
  value = data.edgio_tls_certs.my_certs
}

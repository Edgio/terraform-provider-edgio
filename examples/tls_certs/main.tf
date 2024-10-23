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

provider "edgio" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

data "edgio_tls_certs" "my_certs" {
   environment_id = "enter environment id"
   item_count = 100
}

output "all_my_certs" {
  value = data.edgio_tls_certs.my_certs
}

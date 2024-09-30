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

resource "edgio_purge_cache" "my_purge_cache" {
  environment_id = var.environment_id
  purge_type     = "all_entries"
  values         =  []
}

output "added_purge_cache" {
  value = {
    percentage: edgio_purge_cache.my_purge_cache.progress_percentage
  }
}
 
 
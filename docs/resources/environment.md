---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edgio_environment Resource - terraform-provider-edgio"
subcategory: ""
description: |-
  
---

# edgio_environment (Resource)

Use the `edgio_environment` resource to:
* Create a new environment.
* Update an existing environment.
* Delete an environment.

Learn more about the environment resource in the [Edgio API documentation](https://docs.edg.io/applications/v7/basics/environments).

## Example Usage

```terraform
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
variable "organization_id" { type = string }

provider "edgio" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

resource "edgio_property" "my_property" {
  organization_id = var.organization_id
  slug            = "edgio-environment-example"
}

resource "edgio_environment" "my_env" {
  property_id                 = edgio_property.my_property.id
  name                        = "main"
  only_maintainers_can_deploy = true  
  http_request_logging        = true
}

output "added_env" {
  value = edgio_environment.my_env
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `property_id` (String)

### Optional

- `http_request_logging` (Boolean)
- `only_maintainers_can_deploy` (Boolean)

### Read-Only

- `created_at` (String)
- `default_domain_name` (String)
- `dns_domain_name` (String)
- `id` (String) The ID of this resource.
- `legacy_account_number` (String)
- `pci_compliance` (Boolean)
- `updated_at` (String)


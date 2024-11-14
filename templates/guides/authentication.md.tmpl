---
page_title: "Authentication"
---

# Authentication
The Edgio Terraform provider supports authentication via the Edgio REST API using OAuth 2.0 client credentials. To use the Edgio Terraform provider, you need to define a client ID and client secret in the Edgio Platform. For more information, see the [Edgio API documentation](https://docs.edg.io/applications/v7/rest_api/authentication#~(q~'API*20Clients)).

After creating an API client, you can use the client ID and client secret to initialize the Edgio provider. The provider will automatically handle the OAuth 2.0 flow to obtain an access token and use it for authenticating requests to the Edgio API.


### Requirements
For the Edgio provider to work correctly, the following scopes must be enabled:

- `app.accounts`
- `app.config`


### Resources
The Edgio provider supports the following resources:
* edgio_property
* edgio_environment
* edgio_cdn_configuration
* edgio_tls_cert


### Example

```
terraform {
  required_providers {
    edgio = {
      source = "Edgio/edgio"
      version = "1.0.0"
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
  slug = "edgio-property-example"
}
```

#### Note
In the example above, the `client_id`, `client_secret` amd `organizaiton_id` are passed as variables. You can also set these values in the provider block directly. However, it is recommended to use variables to avoid exposing sensitive information in your Terraform configuration files. See more on how to pass sensitive data to Terraform in the [input variables](https://developer.hashicorp.com/terraform/language/values/variables) document.


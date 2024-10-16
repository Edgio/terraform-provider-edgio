terraform {
  required_providers {
    edgio = {
      source = "hashicorp.com/edu/edgio"
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
  slug = "edgio-config-example"
}

resource "edgio_environment" "my_env" {
  property_id                 = edgio_property.my_property.id
  name                        = "main"
  only_maintainers_can_deploy = false  
  http_request_logging        = true
}

resource "edgio_cdn_configuration" "my_cdn_configuration" {
    environment_id = edgio_environment.my_env.id
    rules = jsonencode(
    [
        {
            "if": [
            {
                "==": [
                {
                    "request": "path"
                },
                "/:path*"
                ]
            },
            {
                "origin": {
                    "set_origin": "edgio_serverless"
                },
                "headers": {
                    "set_request_headers": {
                        "+x-cloud-functions-hint": "app"
                    }
                }
            }
            ]
        }
    ])	
    origins = [
        {
            name: "origin-1",
            type: "customer_origin",
             balancer: "round_robin",	
            override_host_header: "edgio-terraform-example.com",
            pci_certified_shields: false,			
            hosts: [
                {
                    scheme: "https",
                    weight: 100,
                    use_sni: false,
                    balancer: "round_robin",
                    location: [
                    {
                        port: 443,
                        hostname: "origin.edgio-terraform-example.com"
                    }
                    ],
                    max_pool: 0,
                    dns_max_ttl: 3600,
                    dns_min_ttl: 600,
                    max_hard_pool: 10,
                    dns_preference: "prefv4",
                    override_host_header: "edgio-terraform-example.com",
                    sni_hint_and_strict_san_check: "edgio-terraform-example.com"
                }
            ],
        }
    ]

    hostnames = [{
        hostname             = "cdn.edgio-terraform-example.com"
        default_origin_name  = "origin-1"

        tls = {
            npn                = true
            alpn               = true
            protocols          = "TLSv1.2"
            use_sigalgs        = true
            sni                = true
            sni_strict         = true
            sni_host_match     = true
            client_renegotiation = false
            cipher_list        = "ECDHE-RSA-AES128-GCM-SHA256"
        }
    }]		
}

output "added_config" {
  value = edgio_cdn_configuration.my_cdn_configuration
}
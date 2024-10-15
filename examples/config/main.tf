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

resource "edgio_property" "property_name" {
  
}
resource "edgio_environment" "my_environment" {
  name = "my-environment"
}

resource "edgio_cdn_configuration" "my_cdn_configuration" {
  environment_id = var.environment_id
  property_id = edgio_property.property_name.id
  rules = <<EOF
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
]
EOF
  origins = [
  {
    name: "origin-1",
    hosts: [
      {
        scheme: "",
        weight: 200,
        use_sni: false,
        balancer: "",
        location: [
          {
            port: 443,
            hostname: "origin.example.com"
          }
        ],
        max_pool: 0,
        dns_max_ttl: 3600,
        dns_min_ttl: 600,
        max_hard_pool: 10,
        dns_preference: "ipv4",
        override_host_header: "",
        sni_hint_and_strict_san_check: ""
      }
    ],
    balancer: "round_robin",
    override_host_header: "example.com",
    pci_certified_shields: false
  }
]

   hostnames = [{
    hostname             = "cdn.example.com"  # Required hostname
    default_origin_name  = "origin-1"         # Optional default origin name

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

output "my_cdn_configuration" {
  value = edgio_cdn_configuration.my_cdn_configuration
}
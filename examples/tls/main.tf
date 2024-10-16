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
  slug = "test-edgio-property"
}

resource "edgio_environment" "my_env" {
  property_id = edgio_property.my_property.id
  name        = "main"
  only_maintainers_can_deploy = false  
  http_request_logging = true
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

resource "edgio_tls_cert" "my_cert" {
    depends_on = [edgio_cdn_configuration.my_cdn_configuration]
    environment_id    = edgio_environment.my_env.id
    primary_cert      = "-----BEGIN CERTIFICATE-----\nMIIDGjCCAgKgAwIBAgIBAjANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdSdWJ5\nIENBMB4XDTI0MTAxNTA4MTAzOFoXDTI1MTAxNTA4MTAzOFowXDETMBEGCgmSJomT\n8ixkARkWA29yZzEZMBcGCgmSJomT8ixkARkWCXJ1YnktbGFuZzEqMCgGA1UEAwwh\nZWRnaW8tdGVycmFmb3JtLXByb3ZpZGVyLXRlc3QuY29tMIIBIjANBgkqhkiG9w0B\nAQEFAAOCAQ8AMIIBCgKCAQEAxDVYjQfoSNT7QFI7LgFl4Z3R2ye4Co0wnyq5KebD\nRveVOuMIsiTwm//7MrdJ2kYVFRLlFiC9e71OMSRpiBWYgW0l4Zew/L5uI0pLq2zQ\nMpZHsUhJkzVm6H8p59hosNs6G435m561yuGzSY6ze5c0kRWG1DKA2ckqFr8DbHdX\nvH/ao3Gr+XiA1OXQFf4BUJG5pAmMvTOZAGXHcr/4t+aFDoS/V7ZifP7ZvchnntTA\n2aGgpo4Kmxjfx1Kr5jsZkDpcqoltW+NbUJMfNpy/IWdJPQI6Sfzm1ffPTUCzyq7o\nUVpC4aAxDFxileoQ4LsZfkvIZkRO5C6PfoE9X+xt/NlnmwIDAQABozEwLzAOBgNV\nHQ8BAf8EBAMCB4AwHQYDVR0OBBYEFG5oVK9F5TDjWCrtA3sPmgzHUtUOMA0GCSqG\nSIb3DQEBCwUAA4IBAQCav70f8mfvIhoMmaL0NPGUoVadt0E9DEspsA3y5CtO+yHm\nTPJSDpIZlAi2UF4hOuvhR7RIxCeUZVd7C2w3kTkQTttxOLY3D7wVNmnVgTpBzrR5\njhJWg0lsvMC/Kld10v7TzIjEReqj0NQzuJta0etqhBxTDV9e7k5Sk210M0MtFMvV\n2jBWSOOXV8m4EBNN9ooPc5c4rFTxyJTwAil9Iji3VeJ+NvwvW5rLDpemVcfsKCOz\nfWzGofvEfCnxz1sIJPTwFNAY5LtsQsxjk4/FFA6k/x81Ji19Y3Aq/+FgdoQcmBtF\nD4lJhbRqaIp6scJJjiTmintj7y4qbLZZOkJ0nAYx\n-----END CERTIFICATE-----\n"
    intermediate_cert = "-----BEGIN CERTIFICATE-----\nMIIDAjCCAeqgAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdSdWJ5\nIENBMB4XDTE0MTAxNTA4MTAzOFoXDTM0MTAxNTA4MTAzOFowEjEQMA4GA1UEAwwH\nUnVieSBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALz2S/0OVLyn\nJsPDbrAEQkKO+aaRpbBbljFy6vJ49sgoX5BFEtrt16DDJcQktHyJEWLEeWsuG//0\nNpd+FC1GwF4Ns9ziPvTfZgSNYzaGuueNLi8qIyMXCo5F4r8dwF5dKSdvO2VuJzuM\n3a9g6iCXidi+Tvy1Ue6S9G026+WyO20ZA1ovf3G5YyoQA9ZegmvPGg+IzWvCQHP9\nmTYb7MyTC6pM1+WfH5hie05KIM8xohcbu+yBtZ3TnkyOYB63ePW3DEeGoTHoNpj8\nvFOTF44WKljMQhoYxte/n+WWmYoEUd/kV10V6Z6m9uMaagL9NuxWryQuwpd8qTw6\nBBCuSmtJlw0CAwEAAaNjMGEwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC\nAQYwHQYDVR0OBBYEFAw1s247NtKd26sne+81Fu7g/fMFMB8GA1UdIwQYMBaAFAw1\ns247NtKd26sne+81Fu7g/fMFMA0GCSqGSIb3DQEBCwUAA4IBAQB2AoShX2rwzWzY\n8HUYkmdpEfBfKetsWfVrrHnL5hIuwWzzpWY9MHE7VhFUKXSy3RNY0yaWdpozOpz5\nYm77+oJuZgdOyuGLFvySTtElJmK7r8vya5oFGpAgRObzfunL/7wXbXrpTeQFvTZ6\nlGz8fM8GpNIq0v7pxpT+CDBFveYZLL3n0KzuHC4D/YIc4FzeabFGUk9zx1hGLbiD\nh9qHNktax1U6ZXFy8SkEzdXjrf6XbbdY190W7bklNlLaM04Jl9rsGyGgLOfkXIeR\ntOo8WRZeKIjp+PmpUdKUzgo0GRNmVIbax/drc8tTXxXm5NSeo6plFgxG8qpfAGp3\nFmf9Ptew\n-----END CERTIFICATE-----\n"
    private_key       = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxDVYjQfoSNT7QFI7LgFl4Z3R2ye4Co0wnyq5KebDRveVOuMI\nsiTwm//7MrdJ2kYVFRLlFiC9e71OMSRpiBWYgW0l4Zew/L5uI0pLq2zQMpZHsUhJ\nkzVm6H8p59hosNs6G435m561yuGzSY6ze5c0kRWG1DKA2ckqFr8DbHdXvH/ao3Gr\n+XiA1OXQFf4BUJG5pAmMvTOZAGXHcr/4t+aFDoS/V7ZifP7ZvchnntTA2aGgpo4K\nmxjfx1Kr5jsZkDpcqoltW+NbUJMfNpy/IWdJPQI6Sfzm1ffPTUCzyq7oUVpC4aAx\nDFxileoQ4LsZfkvIZkRO5C6PfoE9X+xt/NlnmwIDAQABAoIBAAO8HjjlByNnxnaV\neiHojedrCSUaTvMId/33orms9leh+9m4m6BEer4Fc+MlwQaiIeGaT/kJW4IA+v2N\ne2LHQnVoPfna2NgeydrrHaCgPCBSYv/5Z8khEZnoXcRXhrqjGaqPm8o+DajUfgSu\n7jSyjqIaXkwov/IlVaNENIz6gpWIcvjidTEJ6KJLbC9N5BBs/kuyJobSYd3fJZ+Z\nWXSbVx22SgHyCJgtlE/jhHyeKhQy7MEfuCk6gBiTxcpCLTEOQWMQDbq1wZs4nA3o\n+0rPQ/KgHmm27vZfbFU3i+V+AWSwlXFKd7FWZus2FsJfzYZL3BS1gbFxfeAkfc46\nQ78Y3SkCgYEA4/eLVfvijfmQ7R6iKfkPCwx8bEvQD2YZcBjHzzWY76YQfE/4Ixqa\n2ARKnpmOpRu4oUxioXmgrTUdlC7sVnIpx5s2HPVrLqwoyFXS8Wp0VjOTSf+TGDCy\nBon6Ljjgn6Kj1NhlWsf+Hijt4pZ9e7++LESkCNKIVmCM4Qj5boMJCtMCgYEA3FYI\nlyVF+9gjo+uLK09+1xUMagQ8TbtVk6IA+ZnxtMtikZPkt2OWFIIEU4KGF/k3FvEB\nLnOFQ7C1RhZTsozjpyfcO2oqijkLyI+EgOf7g/ATgm9HtLnrcw8F23EyjX2RxcoL\nIQDv6gHdGPSnuxGv57y/R9LRJjXJRvPrt+LBoxkCgYEA4LfVb1YUNzXrKgNHga6U\nqKSPRkXZfER+EOUsmdLQxnPhzlkaVqhUOVrJn9vpJFLWRpJAq8J0pCk21isHKBPz\noWMcDaHTHTfyH8GSZg41TgAbUheQjYj7BL0glE3XByXQ7/C8wKdilaJtFS6Z1dHm\nikbDmDrI0LTuSqqJDuo2kKcCgYEAlxopKf5V0DCZwIB4IGuUAMxehxYAhQ5D0cr4\nADSineoc3tkdsOaKteW0MdEBRM+UCBefR8vRSGqW5knJfFlChg+/6L8WDVSx0Akc\nRYrR4dlyh7Do6/fUkENtMOCgWogSwCGfIDMUVNaSWdrubEvk5nd6djcNV7brIc2F\nicXoJYECgYBQXP5DYzGWH0V5a/V/W26KNAeB1zdw7F2l975nGkgmIw5QdAsZzWh+\nqZaOiGwO7tml63wiOCF/mJ2jhhio1p+AWfQco9JWwRUQP8p9RSqHlG60fVFIqzGK\noKt8Kxv8gR0m/10zxkFexog/N+jIfr4eMveF0XoiOg455adPAZN2lA==\n-----END RSA PRIVATE KEY-----\n"
}

output "added_property" {
  value = edgio_property.my_property
}
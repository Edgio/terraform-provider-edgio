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

resource "edgio_cdn_configuration" "my_cdn_configuration" {
  environment_id = var.environment_id
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
  },
  {
    "if": [
      {
        "in": [
          {
            "request": "path"
          },
          [
            "/data.json",
            "/favicon.ico"
          ]
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/:path*:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/aleksandar-perkucin-nuxt-app-2/.output/public/:path*:optionalQuery"
            }
          ]
        },
        "origin": {
          "set_origin": "edgio_static"
        },
        "caching": {
          "max_age": {
            "200": "315360000s"
          }
        }
      }
    ]
  },
  {
    "if": [
      {
        "==": [
          {
            "request": "path"
          },
          "/_nuxt/:path*"
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/_nuxt/:path*:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/permanent-env-93628e5b-49b8-4750-9953-f37e3a3e91a7/.output/public/_nuxt/:path*:optionalQuery"
            }
          ]
        },
        "origin": {
          "set_origin": "edgio_permanent_static"
        },
        "caching": {
          "max_age": {
            "200": "1y"
          },
          "client_max_age": "315360000s"
        }
      }
    ]
  },
  {
    "if": [
      {
        "==": [
          {
            "request": "path"
          },
          "/service-worker.js"
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/service-worker.js:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/aleksandar-perkucin-nuxt-app-2/.output/public/_nuxt/service-worker.js:optionalQuery"
            }
          ]
        },
        "origin": {
          "set_origin": "edgio_static"
        },
        "caching": {
          "max_age": {
            "200": "1y"
          },
          "bypass_client_cache": true
        }
      }
    ]
  },
  {
    "if": [
      {
        "==": [
          {
            "request": "path"
          },
          "/__edgio__/devtools/:path*"
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/__edgio__/devtools/:path*:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/aleksandar-perkucin-nuxt-app-2/node_modules/@edgio/devtools/widget/:path*:optionalQuery"
            }
          ]
        },
        "origin": {
          "set_origin": "edgio_static"
        },
        "caching": {
          "max_age": {
            "200": "1y"
          },
          "bypass_client_cache": true
        }
      }
    ]
  },
  {
    "if": [
      {
        "and": [
          {
            "==": [
              {
                "request": "path"
              },
              "/__edgio__/devtools/enable"
            ]
          },
          {
            "===": [
              {
                "request": "method"
              },
              "GET"
            ]
          }
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/:path*:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/:path*:optionalSlash:optionalQuery"
            }
          ],
          "url_redirect": {
            "code": 302,
            "source": "/__edgio__/devtools/enable:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
            "syntax": "path-to-regexp",
            "destination": "/:optionalSlash"
          }
        },
        "caching": {
          "bypass_cache": true,
          "bypass_client_cache": true
        },
        "headers": {
          "add_response_headers": {
            "set-cookie": "edgio_devtools_enabled=true;│ Path=/"
          }
        }
      }
    ]
  },
  {
    "if": [
      {
        "and": [
          {
            "==": [
              {
                "request": "path"
              },
              "/__edgio__/devtools/disable"
            ]
          },
          {
            "===": [
              {
                "request": "method"
              },
              "GET"
            ]
          }
        ]
      },
      {
        "url": {
          "url_rewrite": [
            {
              "source": "/:path*:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
              "syntax": "path-to-regexp",
              "destination": "/:path*:optionalSlash:optionalQuery"
            }
          ],
          "url_redirect": {
            "code": 302,
            "source": "/__edgio__/devtools/disable:optionalSlash(\\/?)?:optionalQuery(\\?.*)?",
            "syntax": "path-to-regexp",
            "destination": "/:optionalSlash"
          }
        },
        "caching": {
          "bypass_cache": true,
          "bypass_client_cache": true
        },
        "headers": {
          "add_response_headers": {
            "set-cookie": "edgio_devtools_enabled=false;│ Path=/"
          }
        }
      }
    ]
  },
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
        "headers": {
          "debug_header": true
        }
      }
    ]
  },
  {
    "if": [
      {
        "=~": [
          {
            "request.header": "host"
          },
          "\\.edgio\\.link|\\.edgio-perma\\.link"
        ]
      },
      {
        "headers": {
          "add_response_headers": {
            "x-robots-tag": "nofollow,│ noindex"
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
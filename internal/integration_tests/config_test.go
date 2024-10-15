package integration_tests

import (
	"fmt"
	"os"
	"terraform-provider-edgio/internal/edgio_provider"
	"terraform-provider-edgio/internal/edgio_provider/utility"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfig_Lifecycle(t *testing.T) {
	client_id := os.Getenv("EDGIO_CLIENT_ID")
	client_secret := os.Getenv("EDGIO_CLIENT_SECRET")
	organization_id := os.Getenv("EDGIO_ORGANIZATION_ID")
	slug_name := fmt.Sprintf("prop-%s", utility.RandomString(10))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: getConfigConfig(client_id, client_secret, organization_id, slug_name, "example-environment", true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"edgio_environment.env_test", "property_id",
						"edgio_property.prop_env_test", "id",
					),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "name", "example-environment"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "can_members_deploy", "true"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "http_request_logging", "true"),
				),
			},
			{
				Config: getConfigConfig(client_id, client_secret, organization_id, slug_name, "updated-environment", false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"edgio_environment.env_test", "property_id",
						"edgio_property.prop_env_test", "id",
					),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "name", "updated-environment"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "can_members_deploy", "false"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "http_request_logging", "false"),
				),
			},
		},
	})
}

func getConfigConfig(clientID, clientSecret, organizationID, propertySlug, envName string, canMembersDeploy, httpRequestLogging bool) string {
	return fmt.Sprintf(`
	provider "edgio" {
		client_id     = "%s"
		client_secret = "%s"
	}

	resource "edgio_property" "prop_env_test" {
		organization_id = "%s"
		slug            = "%s"
	}
		
	resource "edgio_environment" "env_test" {
		property_id         = edgio_property.prop_env_test.id
		name                = "%s"
		can_members_deploy  = %t
		http_request_logging = %t
	}
	
	resource "edgio_cdn_configuration" "my_cdn_configuration" {
  		environment_id = edgio_environment.env_test.id
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
				hosts: [
					{
						scheme: "https",
						weight: 100,
						use_sni: false,
						balancer: "round_robin",
						location: [
						{
							port: 443,
							hostname: "origin.edgio-terraform-provider-test.com"
						}
						],
						max_pool: 0,
						dns_max_ttl: 3600,
						dns_min_ttl: 600,
						max_hard_pool: 10,
						dns_preference: "prefv4",
						override_host_header: "edgio-terraform-provider-test.com",
						sni_hint_and_strict_san_check: "edgio-terraform-provider-test.com"
					}
				],
				balancer: "round_robin",
				override_host_header: "edgio-terraform-provider-test.com",
				pci_certified_shields: false
			}
		]

		hostnames = [{
			hostname             = "cdn.edgio-terraform-provider-test.com"  # Required hostname
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
	}`, clientID, clientSecret, organizationID, propertySlug, envName, canMembersDeploy, httpRequestLogging)
}

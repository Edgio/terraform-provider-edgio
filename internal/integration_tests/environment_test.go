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

func TestEnvironment_Lifecycle(t *testing.T) {
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
				Config: getEnvConfig(client_id, client_secret, organization_id, slug_name, "example-environment", true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"edgio_environment.env_test", "property_id",
						"edgio_property.prop_env_test", "id",
					),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "name", "example-environment"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "only_maintainers_can_deploy", "true"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "http_request_logging", "true"),
				),
			},
			{
				Config: getEnvConfig(client_id, client_secret, organization_id, slug_name, "updated-environment", false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"edgio_environment.env_test", "property_id",
						"edgio_property.prop_env_test", "id",
					),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "name", "updated-environment"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "only_maintainers_can_deploy", "false"),
					resource.TestCheckResourceAttr("edgio_environment.env_test", "http_request_logging", "false"),
				),
			},
		},
	})
}

func getEnvConfig(clientID, clientSecret, organizationID, propertySlug, envName string, canMembersDeploy, httpRequestLogging bool) string {
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
		only_maintainers_can_deploy  = %t
		http_request_logging = %t
	}`, clientID, clientSecret, organizationID, propertySlug, envName, canMembersDeploy, httpRequestLogging)
}

package integration_tests

import (
	"fmt"
	"os"
	"testing"

	"terraform-provider-edgio/internal/edgio_provider"
	"terraform-provider-edgio/internal/edgio_provider/utility"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProperty_Lifecycle(t *testing.T) {
	client_id := os.Getenv("EDGIO_CLIENT_ID")
	client_secret := os.Getenv("EDGIO_CLIENT_SECRET")
	organization_id := os.Getenv("EDGIO_ORGANIZATION_ID")
	slug_name := fmt.Sprintf("prop-%s", utility.RandomString(10))
	updated_slug_name := fmt.Sprintf("%s-updated", slug_name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: getPropertyConfig(client_id, client_secret, organization_id, slug_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_property.prop_test", "organization_id", organization_id),
					resource.TestCheckResourceAttr("edgio_property.prop_test", "slug", slug_name),
				),
			},
			{
				ResourceName:      "edgio_property.prop_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: getPropertyConfig(client_id, client_secret, organization_id, updated_slug_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_property.prop_test", "organization_id", organization_id),
					resource.TestCheckResourceAttr("edgio_property.prop_test", "slug", updated_slug_name),
				),
			},
		},
	})
}

func getPropertyConfig(clientID, clientSecret, organizationID, propertySlug string) string {
	return fmt.Sprintf(`
	provider "edgio" {
		client_id     = "%s"
		client_secret = "%s"
	}

	resource "edgio_property" "prop_test" {
		organization_id = "%s"
		slug            = "%s"		
	}`, clientID, clientSecret, organizationID, propertySlug)
}

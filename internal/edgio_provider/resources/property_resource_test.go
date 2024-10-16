package resources_test

import (
	"testing"
	"time"

	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider"
	"terraform-provider-edgio/internal/edgio_provider/utility"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func mockAllPropertyMethods(mockClient *edgio_api.MockEdgioClient, methods ...utility.MockMethod) {
	fixedTime := time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC)
	updatedAt := fixedTime
	property := &dtos.Property{
		OrganizationID: "org-123",
		Slug:           "example-slug",
		Id:             "property-123",
		IdLink:         "property-link-123",
		CreatedAt:      fixedTime,
		UpdatedAt:      updatedAt,
	}

	for _, method := range methods {
		switch method {
		case utility.MockCreate:
			mockClient.On("CreateProperty", mock.Anything, "org-123", "example-slug").Return(property, nil)
		case utility.MockGet:
			mockClient.On("GetProperty", mock.Anything, "property-123").Return(property, nil)
		case utility.MockUpdate:
			mockClient.On("UpdateProperty", mock.Anything, "property-123", "new-slug").Run(func(args mock.Arguments) {
				property.Slug = "new-slug"
				property.UpdatedAt = time.Now()
			}).Return(property, nil)
		case utility.MockDelete:
			mockClient.On("DeleteProperty", "property-123").Return(nil)
		}
	}
}

func TestPropertyResource_Lifecycle(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	mockAllPropertyMethods(mockClient, utility.MockCreate, utility.MockGet, utility.MockUpdate, utility.MockDelete)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.NewMockedProvider(mockClient)),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

                resource "edgio_property" "test" {
                    organization_id = "org-123"
                    slug            = "example-slug"
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_property.test", "organization_id", "org-123"),
					resource.TestCheckResourceAttr("edgio_property.test", "slug", "example-slug"),
				),
			},
			{
				ResourceName:      "edgio_property.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

                resource "edgio_property" "test" {
                    organization_id = "org-123"
                    slug            = "new-slug"
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_property.test", "organization_id", "org-123"),
					resource.TestCheckResourceAttr("edgio_property.test", "slug", "new-slug"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

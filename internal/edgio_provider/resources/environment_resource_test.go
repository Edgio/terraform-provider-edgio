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

func mockAllEnvironmentMethods(mockClient *edgio_api.MockEdgioClient, methods ...utility.MockMethod) {
	fixedTime := time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC)
	updatedAt := fixedTime
	environment := &dtos.Environment{
		Id:                       "env-123",
		PropertyID:               "property-123",
		LegacyAccountNumber:      "legacy-123",
		Name:                     "example-environment",
		CanMembersDeploy:         true,
		OnlyMaintainersCanDeploy: false,
		HttpRequestLogging:       true,
		DefaultDomainName:        "example.com",
		PciCompliance:            false,
		DnsDomainName:            "example.dns.com",
		CreatedAt:                fixedTime,
		UpdatedAt:                updatedAt,
	}

	for _, method := range methods {
		switch method {
		case utility.MockCreate:
			mockClient.On("CreateEnvironment", "property-123", "example-environment", false, true).Return(environment, nil)
		case utility.MockGet:
			mockClient.On("GetEnvironment", "env-123").Return(environment, nil)
		case utility.MockUpdate:
			mockClient.On("UpdateEnvironment", "env-123", "updated-environment", false, false, false).Run(func(args mock.Arguments) {
				environment.Name = "updated-environment"
				environment.CanMembersDeploy = false
				environment.HttpRequestLogging = false
				environment.UpdatedAt = time.Now()
			}).Return(environment, nil)
		case utility.MockDelete:
			mockClient.On("DeleteEnvironment", "env-123").Return(nil)
		}
	}
}

func TestEnvironmentResource_Lifecycle(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	mockAllEnvironmentMethods(mockClient, utility.MockCreate, utility.MockGet, utility.MockUpdate, utility.MockDelete)

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

				resource "edgio_environment" "test" {
					property_id         = "property-123"
					name                = "example-environment"
					only_maintainers_can_deploy = false
					http_request_logging = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_environment.test", "property_id", "property-123"),
					resource.TestCheckResourceAttr("edgio_environment.test", "name", "example-environment"),
					resource.TestCheckResourceAttr("edgio_environment.test", "only_maintainers_can_deploy", "false"),
					resource.TestCheckResourceAttr("edgio_environment.test", "http_request_logging", "true"),
				),
			},
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_environment" "test" {
					property_id         = "property-123"
					name                = "updated-environment"
					only_maintainers_can_deploy  = false
					http_request_logging = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_environment.test", "property_id", "property-123"),
					resource.TestCheckResourceAttr("edgio_environment.test", "name", "updated-environment"),
					resource.TestCheckResourceAttr("edgio_environment.test", "only_maintainers_can_deploy", "false"),
					resource.TestCheckResourceAttr("edgio_environment.test", "http_request_logging", "false"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

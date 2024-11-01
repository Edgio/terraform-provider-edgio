package resources_test

import (
	"encoding/json"
	"testing"

	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider"
	"terraform-provider-edgio/internal/edgio_provider/utility"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func mockAllCDNConfigurationMethods(mockClient *edgio_api.MockEdgioClient, methods ...utility.MockMethod) {
	cdnConfig := &dtos.CDNConfiguration{
		ConfigurationID: "config-123",
		EnvironmentID:   "env-123",
		Rules: json.RawMessage(`
        {
            "test":123
        }`),
		Origins: []dtos.Origin{
			{
				Name: "origin-1",
				Hosts: []dtos.Host{
					{
						Weight: utility.ToPtr(int64(200)),
						UseSNI: utility.ToPtr(false),
						Location: utility.ToPtr([]dtos.Location{
							{
								Port:     utility.ToPtr(int64(443)),
								Hostname: utility.ToPtr("origin.example.com"),
							},
						}),
						MaxPool:       utility.ToPtr(int64(0)),
						DNSMaxTTL:     utility.ToPtr(int64(3600)),
						DNSMinTTL:     utility.ToPtr(int64(600)),
						MaxHardPool:   utility.ToPtr(int64(10)),
						DNSPreference: utility.ToPtr("ipv4"),
					},
				},
				Balancer:            utility.ToPtr("round_robin"),
				OverrideHostHeader:  utility.ToPtr("example.com"),
				PciCertifiedShields: utility.ToPtr(false),
			},
		},
		Hostnames: []dtos.Hostname{
			{
				Hostname:          utility.ToPtr("cdn.example.com"),
				DefaultOriginName: utility.ToPtr("origin-1"),
				TLS: &dtos.TLS{
					NPN:                 utility.ToPtr(true),
					ALPN:                utility.ToPtr(true),
					Protocols:           utility.ToPtr("TLSv1.2"),
					UseSigAlgs:          utility.ToPtr(true),
					SNI:                 utility.ToPtr(true),
					SniStrict:           utility.ToPtr(true),
					SniHostMatch:        utility.ToPtr(true),
					ClientRenegotiation: utility.ToPtr(false),
					CipherList:          utility.ToPtr("ECDHE-RSA-AES128-GCM-SHA256"),
				},
			},
		},
	}

	for _, method := range methods {
		switch method {
		case utility.MockUpload:
			mockClient.On("UploadCdnConfiguration", mock.AnythingOfType("*dtos.CDNConfiguration")).Return(cdnConfig, nil)
		case utility.MockGet:
			mockClient.On("GetCDNConfiguration", "config-123").Return(cdnConfig, nil)
		}
	}
}

func TestCDNConfigurationResource_Lifecycle(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	mockAllCDNConfigurationMethods(mockClient, utility.MockUpload, utility.MockGet)

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

				resource "edgio_cdn_configuration" "test" {			
  					environment_id = "env-123"
					rules = jsonencode({
						"test": 123
					})					
					origins = [
						{
							name: "origin-1",
							hosts: [
							{					
								weight: 200,
								use_sni: false,					
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
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_cdn_configuration.test", "environment_id", "env-123"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

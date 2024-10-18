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
						Scheme:   "",
						Weight:   200,
						UseSNI:   false,
						Balancer: "",
						Location: []dtos.Location{
							{
								Port:     443,
								Hostname: "origin.example.com",
							},
						},
						MaxPool:                  0,
						DNSMaxTTL:                3600,
						DNSMinTTL:                600,
						MaxHardPool:              10,
						DNSPreference:            "ipv4",
						OverrideHostHeader:       "",
						SNIHintAndStrictSanCheck: "",
					},
				},
				Balancer:            utility.PtrString("round_robin"),
				OverrideHostHeader:  utility.PtrString("example.com"),
				PciCertifiedShields: utility.PtrBool(false),
			},
		},
		Hostnames: []dtos.Hostname{
			{
				Hostname:          "cdn.example.com",
				DefaultOriginName: "origin-1",
				TLS: dtos.TLS{
					NPN:                 true,
					ALPN:                true,
					Protocols:           "TLSv1.2",
					UseSigAlgs:          true,
					SNI:                 true,
					SniStrict:           true,
					SniHostMatch:        true,
					ClientRenegotiation: false,
					CipherList:          "ECDHE-RSA-AES128-GCM-SHA256",
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
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_cdn_configuration.test", "environment_id", "env-123"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

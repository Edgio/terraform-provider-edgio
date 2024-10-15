package resources_test

import (
	"testing"
	"time"

	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestTLSCertsResource_Generate(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	// Mock data
	fixedTime := time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC)
	generatedCert := &dtos.TLSCertResponse{
		ID:               "cert-123",
		EnvironmentID:    "env-123",
		Expiration:       fixedTime.Add(365 * 24 * time.Hour).Format(time.RFC3339),
		Status:           "activated",
		Generated:        true,
		Serial:           "1234567890",
		CommonName:       "example.com",
		AlternativeNames: []string{"www.example.com"},
		CreatedAt:        fixedTime.Format(time.RFC3339),
		UpdatedAt:        fixedTime.Format(time.RFC3339),
	}

	mockClient.On("GenerateTlsCert", "env-123").Return(generatedCert, nil)
	mockClient.On("GetTlsCert", "cert-123").Return(generatedCert, nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.NewMockedProvider(mockClient)),
		},
		Steps: []resource.TestStep{
			// Test generating a TLS cert
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "generated" {
					environment_id = "env-123"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "environment_id", "env-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "id", "cert-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "status", "activated"),
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "generated", "true"),
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "common_name", "example.com"),
				),
			},
			// Test reading the generated cert
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "generated" {
					environment_id = "env-123"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "id", "cert-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.generated", "serial", "1234567890"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

func TestTLSCertsResource_Upload(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	fixedTime := time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC)

	uploadedCert := &dtos.TLSCertResponse{
		ID:               "cert-456",
		EnvironmentID:    "env-123",
		Expiration:       fixedTime.Add(365 * 24 * time.Hour).Format(time.RFC3339),
		Status:           "activated",
		PrimaryCert:      "-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIUJx...\n-----END CERTIFICATE-----",
		IntermediateCert: "-----BEGIN CERTIFICATE-----\nMIIDIjCCAg...\n-----END CERTIFICATE-----",
		Generated:        false,
		Serial:           "0987654321",
		CommonName:       "example.org",
		AlternativeNames: []string{"www.example.org"},
		CreatedAt:        fixedTime.Format(time.RFC3339),
		UpdatedAt:        fixedTime.Format(time.RFC3339),
	}

	mockClient.On("UploadTlsCert", mock.AnythingOfType("dtos.UploadTlsCertRequest")).Return(uploadedCert, nil)
	mockClient.On("GetTlsCert", "cert-456").Return(uploadedCert, nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.NewMockedProvider(mockClient)),
		},
		Steps: []resource.TestStep{
			// Test uploading a TLS cert
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "uploaded" {
					environment_id     = "env-123"
					primary_cert       = "-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIUJx...\n-----END CERTIFICATE-----"
					intermediate_cert  = "-----BEGIN CERTIFICATE-----\nMIIDIjCCAg...\n-----END CERTIFICATE-----"
					private_key        = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADAN...\n-----END PRIVATE KEY-----"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "environment_id", "env-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "id", "cert-456"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "status", "activated"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "generated", "false"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "common_name", "example.org"),
				),
			},
			// Test reading the uploaded cert
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "uploaded" {
					environment_id     = "env-123"
					primary_cert       = "-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIUJx...\n-----END CERTIFICATE-----"
					intermediate_cert  = "-----BEGIN CERTIFICATE-----\nMIIDIjCCAg...\n-----END CERTIFICATE-----"
					private_key        = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADAN...\n-----END PRIVATE KEY-----"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "id", "cert-456"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "serial", "0987654321"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

func TestTLSCertsResource_UploadUpdate(t *testing.T) {
	mockClient := new(edgio_api.MockEdgioClient)

	fixedTime := time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC)

	uploadedCert := &dtos.TLSCertResponse{
		ID:               "cert-456",
		EnvironmentID:    "env-123",
		Expiration:       fixedTime.Add(365 * 24 * time.Hour).Format(time.RFC3339),
		Status:           "activated",
		PrimaryCert:      "-----BEGIN CERTIFICATE-----1",
		IntermediateCert: "-----BEGIN CERTIFICATE-----2",
		Generated:        false,
		Serial:           "0987654321",
		CommonName:       "example.org",
		AlternativeNames: []string{"www.example.org"},
		CreatedAt:        fixedTime.Format(time.RFC3339),
		UpdatedAt:        fixedTime.Format(time.RFC3339),
	}

	changedCert := &dtos.TLSCertResponse{
		ID:               "cert-456",
		EnvironmentID:    "env-123",
		Expiration:       fixedTime.Add(365 * 24 * time.Hour).Format(time.RFC3339),
		Status:           "activated",
		PrimaryCert:      "Changed",
		IntermediateCert: "-----BEGIN CERTIFICATE-----2",
		Generated:        false,
		Serial:           "0987654321",
		CommonName:       "example.org",
		AlternativeNames: []string{"www.example.org"},
		CreatedAt:        fixedTime.Format(time.RFC3339),
		UpdatedAt:        fixedTime.Format(time.RFC3339),
	}

	mockClient.On("UploadTlsCert", mock.AnythingOfType("dtos.UploadTlsCertRequest")).Return(uploadedCert, nil)
	mockClient.On("GetTlsCert", "cert-456").Return(uploadedCert, nil).Once()
	mockClient.On("GetTlsCert", "cert-456").Return(changedCert, nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"edgio": providerserver.NewProtocol6WithError(edgio_provider.NewMockedProvider(mockClient)),
		},
		Steps: []resource.TestStep{
			// Test uploading a TLS cert
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "uploaded" {
					environment_id     = "env-123"
					primary_cert       = "-----BEGIN CERTIFICATE-----1"
					intermediate_cert  = "-----BEGIN CERTIFICATE-----2"
					private_key        = "-----BEGIN PRIVATE KEY-----3"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "environment_id", "env-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "id", "cert-456"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "status", "activated"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "generated", "false"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "common_name", "example.org"),
				),
			},
			{
				Config: `
				provider "edgio" {
					client_id     = "mock-client-id"
					client_secret = "mock-client-secret"
				}

				resource "edgio_tls_cert" "uploaded" {
					environment_id     = "env-123"
					primary_cert  	   = "Changed"
					intermediate_cert  = "-----BEGIN CERTIFICATE-----2"
					private_key        = "-----BEGIN PRIVATE KEY-----3"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "environment_id", "env-123"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "id", "cert-456"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "status", "activated"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "generated", "false"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "common_name", "example.org"),
					resource.TestCheckResourceAttr("edgio_tls_cert.uploaded", "primary_cert", "Changed"),
				),
			},
		},
	})

	mockClient.AssertExpectations(t)
}

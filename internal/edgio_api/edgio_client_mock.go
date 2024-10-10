package edgio_api

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api/dtos"

	"github.com/stretchr/testify/mock"
)

type MockEdgioClient struct {
	mock.Mock
}

func (m *MockEdgioClient) GetProperty(ctx context.Context, propertyID string) (*dtos.Property, error) {
	args := m.Called(ctx, propertyID)
	return args.Get(0).(*dtos.Property), args.Error(1)
}

func (m *MockEdgioClient) GetProperties(page int, pageSize int, organizationID string) (*dtos.Properties, error) {
	args := m.Called(page, pageSize, organizationID)
	return args.Get(0).(*dtos.Properties), args.Error(1)
}

func (m *MockEdgioClient) CreateProperty(ctx context.Context, organizationID, slug string) (*dtos.Property, error) {
	args := m.Called(ctx, organizationID, slug)
	return args.Get(0).(*dtos.Property), args.Error(1)
}

func (m *MockEdgioClient) DeleteProperty(propertyID string) error {
	args := m.Called(propertyID)
	return args.Error(0)
}

func (m *MockEdgioClient) UpdateProperty(ctx context.Context, propertyID string, slug string) (*dtos.Property, error) {
	args := m.Called(ctx, propertyID, slug)
	return args.Get(0).(*dtos.Property), args.Error(1)
}

func (m *MockEdgioClient) GetEnvironments(page, pageSize int, propertyID string) (*dtos.EnvironmentsResponse, error) {
	args := m.Called(page, pageSize, propertyID)
	return args.Get(0).(*dtos.EnvironmentsResponse), args.Error(1)
}

func (m *MockEdgioClient) GetEnvironment(environmentID string) (*dtos.Environment, error) {
	args := m.Called(environmentID)
	return args.Get(0).(*dtos.Environment), args.Error(1)
}

func (m *MockEdgioClient) CreateEnvironment(propertyID, name string, canMembersDeploy, onlyMaintainersCanDeploy, httpRequestLogging bool) (*dtos.Environment, error) {
	args := m.Called(propertyID, name, canMembersDeploy, onlyMaintainersCanDeploy, httpRequestLogging)
	return args.Get(0).(*dtos.Environment), args.Error(1)
}

func (m *MockEdgioClient) UpdateEnvironment(environmentID, name string, canMembersDeploy, httpRequestLogging, preserveCache bool) (*dtos.Environment, error) {
	args := m.Called(environmentID, name, canMembersDeploy, httpRequestLogging, preserveCache)
	return args.Get(0).(*dtos.Environment), args.Error(1)
}

func (m *MockEdgioClient) DeleteEnvironment(environmentID string) error {
	args := m.Called(environmentID)
	return args.Error(0)
}

func (m *MockEdgioClient) PurgeCache(purgeRequest *dtos.PurgeRequest) (*dtos.PurgeResponse, error) {
	args := m.Called(purgeRequest)
	return args.Get(0).(*dtos.PurgeResponse), args.Error(1)
}

func (m *MockEdgioClient) GetPurgeStatus(requestId string) (*dtos.PurgeResponse, error) {
	args := m.Called(requestId)
	return args.Get(0).(*dtos.PurgeResponse), args.Error(1)
}

func (m *MockEdgioClient) GetTlsCert(tlsCertId string) (*dtos.TLSCertResponse, error) {
	args := m.Called(tlsCertId)
	return args.Get(0).(*dtos.TLSCertResponse), args.Error(1)
}

func (m *MockEdgioClient) UploadTlsCert(req dtos.UploadTlsCertRequest) (*dtos.TLSCertResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dtos.TLSCertResponse), args.Error(1)
}

func (m *MockEdgioClient) GenerateTlsCert(environmentId string) (*dtos.TLSCertResponse, error) {
	args := m.Called(environmentId)
	return args.Get(0).(*dtos.TLSCertResponse), args.Error(1)
}

func (m *MockEdgioClient) GetTlsCerts(page int, pageSize int, environmentID string) (*dtos.TLSCertSResponse, error) {
	args := m.Called(page, pageSize, environmentID)
	return args.Get(0).(*dtos.TLSCertSResponse), args.Error(1)
}

func (m *MockEdgioClient) UploadCdnConfiguration(config *dtos.CDNConfiguration) (*dtos.CDNConfiguration, error) {
	args := m.Called(config)
	return args.Get(0).(*dtos.CDNConfiguration), args.Error(1)
}

func (m *MockEdgioClient) GetCDNConfiguration(configID string) (*dtos.CDNConfiguration, error) {
	args := m.Called(configID)
	return args.Get(0).(*dtos.CDNConfiguration), args.Error(1)
}

// Ensure MockEdgioClient implements EdgioClientInterface.
var _ EdgioClientInterface = (*MockEdgioClient)(nil)

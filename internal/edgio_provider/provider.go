package edgio_provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_provider/data_sources"
	"terraform-provider-edgio/internal/edgio_provider/resources"
)

// Provider implements the provider.Provider interface.
type Provider struct {
	client edgio_api.EdgioClientInterface
}

func New() provider.Provider {
	return &Provider{}
}

func NewMockedProvider(client edgio_api.EdgioClientInterface) provider.Provider {
	return &Provider{client: client}
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &Provider{}
)

// Metadata returns the provider's metadata.
func (p *Provider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	// Set the provider metadata (optional)
	response.TypeName = "edgio"
	response.Version = "0.1.0"
}

// Configure configures the provider with user-provided configuration.
func (p *Provider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config struct {
		ClientID     types.String `tfsdk:"client_id"`
		ClientSecret types.String `tfsdk:"client_secret"`
	}

	if err := request.Config.Get(ctx, &config); err != nil {
		response.Diagnostics.AddError("Failed to read provider configuration", "error")
		return
	}

	// For mock we don't need to create a new client, as mock
	// will handle all the calls
	if p.client == nil {
		client := edgio_api.NewEdgioClient(
			config.ClientID.ValueString(),
			config.ClientSecret.ValueString(),
			"https://id.edgio.app/connect/token",
			"https://edgioapis.com",
		)
		p.client = client
	}
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return data_sources.NewPropertiesDataSource(p.client)
		},
		func() datasource.DataSource {
			return data_sources.NewEnvironmentsDataSource(p.client)
		},
		func() datasource.DataSource {
			return data_sources.NewTlsCertsDataSource(p.client)
		},
	}
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return resources.NewPropertyResource(p.client)
		},
		func() resource.Resource {
			return resources.NewEnvironmentResource(p.client)
		},
		func() resource.Resource {
			return resources.NewTLSCertsResource(p.client)
		},
		func() resource.Resource {
			return resources.NewCDNConfigurationResource(p.client)
		},
	}
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for OAuth2 authentication.",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret for OAuth2 authentication.",
				Required:            true,
			},
		},
	}
}

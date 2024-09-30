package edgio_provider

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_provider/data_sources"
	"terraform-provider-edgio/internal/edgio_provider/resources"

	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Provider implements the provider.Provider interface.
type Provider struct {
	client *edgio_api.EdgioClient
}

// New creates a new instance of the provider.
func New() provider.Provider {
	return &Provider{}
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &Provider{}
	_ provider.ProviderWithFunctions = &Provider{}
)

// Metadata returns the provider's metadata.
func (p *Provider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	// Set the provider metadata (optional)
	response.TypeName = "edgio"
}

// Configure configures the provider with user-provided configuration.
func (p *Provider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config struct {
		ClientID     types.String `tfsdk:"client_id"`
		ClientSecret types.String `tfsdk:"client_secret"`
	}

	// Read configuration
	if err := request.Config.Get(ctx, &config); err != nil {
		response.Diagnostics.AddError("Failed to read provider configuration", "error")
		return
	}

	// Create a new OAuthClient
	client := edgio_api.NewEdgioClient(
		config.ClientID.ValueString(),
		config.ClientSecret.ValueString(),
		"https://id.edgio.app/connect/token",
		"https://edgioapis.com",
	)

	p.client = client
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
			return resources.NewPurgeCacheResource(p.client)
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

func (p *Provider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function {
			return NewComputeTaxFunction()
		},
	}
}

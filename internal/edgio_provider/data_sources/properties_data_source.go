package data_sources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-edgio/internal/edgio_provider/models"
)

type PropertiesDataSource struct {
	client *edgio_api.EdgioClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &PropertiesDataSource{}
)

func NewPropertiesDataSource(client *edgio_api.EdgioClient) *PropertiesDataSource {
	return &PropertiesDataSource{
		client: client,
	}
}

func (d *PropertiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "edgio_properties"
}

func (d *PropertiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:    true,
				Description: "The organization ID to fetch properties for.",
			},
			"properties": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"slug": schema.StringAttribute{
							Computed: true,
						},
						"property_id": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"links": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"first": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{
											Computed: true,
										},
										"description": schema.StringAttribute{
											Computed: true,
										},
										"base_path": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"next": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{
											Computed: true,
										},
										"description": schema.StringAttribute{
											Computed: true,
										},
										"base_path": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"previous": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{
											Computed: true,
										},
										"description": schema.StringAttribute{
											Computed: true,
										},
										"base_path": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"last": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{
											Computed: true,
										},
										"description": schema.StringAttribute{
											Computed: true,
										},
										"base_path": schema.StringAttribute{
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *PropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.PropertiesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	organizationID := state.OrganizationID.ValueString()

	properties, err := d.client.GetProperties(1, 10, organizationID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading properties", err.Error())
		return
	}

	// Map response body to model
	for _, property := range properties.Items {
		propertyState := models.PropertyModel{
			ID:             types.StringValue(property.ID),
			Slug:           types.StringValue(property.Slug),
			PropertyID:     types.StringValue(property.PropertyID),
			OrganizationID: types.StringValue(property.OrganizationID),
			Type:           types.StringValue(property.Type),
			CreatedAt:      types.StringValue(property.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:      types.StringValue(property.UpdatedAt.Format(time.RFC3339)),
		}

		state.Properties = append(state.Properties, propertyState)
	}

	state.Links = models.PropertiesLinksModel{
		First: models.PropertiesLinkModel{
			Href:        types.StringValue(properties.Links.First.Href),
			Description: types.StringValue(properties.Links.First.Description),
			BasePath:    types.StringValue(properties.Links.First.BasePath),
		},
		Next: models.PropertiesLinkModel{
			Href:        types.StringValue(properties.Links.Next.Href),
			Description: types.StringValue(properties.Links.Next.Description),
			BasePath:    types.StringValue(properties.Links.Next.BasePath),
		},
		Previous: models.PropertiesLinkModel{
			Href:        types.StringValue(properties.Links.Previous.Href),
			Description: types.StringValue(properties.Links.Previous.Description),
			BasePath:    types.StringValue(properties.Links.Previous.BasePath),
		},
		Last: models.PropertiesLinkModel{
			Href:        types.StringValue(properties.Links.Last.Href),
			Description: types.StringValue(properties.Links.Last.Description),
			BasePath:    types.StringValue(properties.Links.Last.BasePath),
		},
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

package data_sources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"time"

	"terraform-provider-edgio/internal/edgio_provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Required: true,
				Description: `An organization's system-defined ID (e.g., 12345678-1234-1234-1234-1234567890ab).
					 From the Edgio Console, navigate to the desired organization and then click Settings. 
					 It is listed under Organization ID."`,
			},
			"item_count": schema.Int32Attribute{
				Computed:    true,
				Description: `The total number of items in the list.`,
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
			"properties": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The resource's system-defined ID.",
						},
						"id_link": schema.StringAttribute{
							Computed:    true,
							Description: "The resource's relative path.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The resource's type.",
						},
						"slug": schema.StringAttribute{
							Computed:    true,
							Description: "The property's name.",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: `An organization's system-defined ID (e.g., 12345678-1234-1234-1234-1234567890ab).
							From the Edgio Console, navigate to the desired organization and then click Settings. 
							It is listed under Organization ID.`,
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The property's creation date and time (UTC).",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The property's last modification date and time (UTC).",
						},
					},
				},
			},
		},
	}
}

func (d *PropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var organizationID string
	diags := req.Config.GetAttribute(ctx, path.Root("organization_id"), &organizationID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := d.client.GetProperties(1, 10, organizationID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading properties", err.Error())
		return
	}

	state := models.PropertiesModel{
		OrganizationID: types.StringValue(organizationID),
		ItemCount:      types.Int32Value(int32(properties.TotalItems)),
		Properties:     []models.PropertyModel{},
		Links: models.PropertiesLinksModel{
			First:    models.PropertiesLinkModel{},
			Next:     models.PropertiesLinkModel{},
			Previous: models.PropertiesLinkModel{},
			Last:     models.PropertiesLinkModel{},
		},
	}

	for _, property := range properties.Items {
		propertyState := models.PropertyModel{
			Id:             types.StringValue(property.Id),
			IdLink:         types.StringValue(property.IdLink),
			Slug:           types.StringValue(property.Slug),
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

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

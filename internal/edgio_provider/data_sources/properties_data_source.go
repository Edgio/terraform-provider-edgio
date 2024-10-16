package data_sources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"time"

	"terraform-provider-edgio/internal/edgio_provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PropertiesDataSource struct {
	client edgio_api.EdgioClientInterface
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &PropertiesDataSource{}
)

func NewPropertiesDataSource(client edgio_api.EdgioClientInterface) *PropertiesDataSource {
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
				Required:    true,
				Description: `The total number of items to load.`,
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
	var state models.PropertiesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := d.client.GetProperties(0, int(state.ItemCount.ValueInt32()), state.OrganizationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading properties", err.Error())
		return
	}

	newState := models.PropertiesModel{
		OrganizationID: state.OrganizationID,
		ItemCount:      state.ItemCount,
		Properties:     []models.PropertyModel{},
	}

	for _, property := range properties.Items {
		propertyState := models.PropertyModel{
			Id:             types.StringValue(property.Id),
			IdLink:         types.StringValue(property.IdLink),
			Slug:           types.StringValue(property.Slug),
			OrganizationID: types.StringValue(property.OrganizationID),
			CreatedAt:      types.StringValue(property.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:      types.StringValue(property.UpdatedAt.Format(time.RFC3339)),
		}

		newState.Properties = append(newState.Properties, propertyState)
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

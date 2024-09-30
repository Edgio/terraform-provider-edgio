package data_sources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvironmentsDataSource struct {
	client *edgio_api.EdgioClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &EnvironmentsDataSource{}
)

func NewEnvironmentsDataSource(client *edgio_api.EdgioClient) *EnvironmentsDataSource {
	return &EnvironmentsDataSource{
		client: client,
	}
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "edgio_environments"
}

func (d *EnvironmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required:    true,
				Description: `The ID of the property to filter environments by.`,
			},
			"item_count": schema.Int32Attribute{
				Computed:    true,
				Description: `The total number of environments.`,
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
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The environment's system-defined ID.",
						},
						"property_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the property associated with the environment.",
						},
						"legacy_account_number": schema.StringAttribute{
							Computed:    true,
							Description: "The legacy account number for the environment.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the environment.",
						},
						"can_members_deploy": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if members can deploy to the environment.",
						},
						"only_maintainers_can_deploy": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if only maintainers can deploy to the environment.",
						},
						"http_request_logging": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if HTTP request logging is enabled for the environment.",
						},
						"default_domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "The default domain name for the environment.",
						},
						"pci_compliance": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if the environment is PCI compliant.",
						},
						"dns_domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "The DNS domain name for the environment.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The environment's creation date and time (UTC).",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The environment's last modification date and time (UTC).",
						},
					},
				},
			},
		},
	}
}

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var propertyID string
	diags := req.Config.GetAttribute(ctx, path.Root("property_id"), &propertyID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := d.client.GetEnvironments(1, 10, propertyID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading environments", err.Error())
		return
	}

	state := models.EnvironmentsModel{
		Id:           types.StringValue(propertyID),
		TotalItems:   types.Int32Value(int32(environments.TotalItems)),
		Environments: []models.EnvironmentModel{},
		Links: models.EnvironmentsLinksModel{
			First:    models.EnvironmentsLinkModel{},
			Next:     models.EnvironmentsLinkModel{},
			Previous: models.EnvironmentsLinkModel{},
			Last:     models.EnvironmentsLinkModel{},
		},
	}

	for _, environment := range environments.Items {
		envState := models.EnvironmentModel{
			Id:                       types.StringValue(environment.Id),
			PropertyID:               types.StringValue(environment.PropertyID),
			LegacyAccountNumber:      types.StringValue(environment.LegacyAccountNumber),
			Name:                     types.StringValue(environment.Name),
			CanMembersDeploy:         types.BoolValue(environment.CanMembersDeploy),
			OnlyMaintainersCanDeploy: types.BoolValue(environment.OnlyMaintainersCanDeploy),
			HttpRequestLogging:       types.BoolValue(environment.HttpRequestLogging),
			DefaultDomainName:        types.StringValue(environment.DefaultDomainName),
			PciCompliance:            types.BoolValue(environment.PciCompliance),
			DnsDomainName:            types.StringValue(environment.DnsDomainName),
			CreatedAt:                types.StringValue(environment.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:                types.StringValue(environment.UpdatedAt.Format(time.RFC3339)),
		}

		state.Environments = append(state.Environments, envState)
	}

	state.Links = models.EnvironmentsLinksModel{
		First: models.EnvironmentsLinkModel{
			Href:        types.StringValue(environments.Links.First.Href),
			Description: types.StringValue(environments.Links.First.Description),
			BasePath:    types.StringValue(environments.Links.First.BasePath),
		},
		Next: models.EnvironmentsLinkModel{
			Href:        types.StringValue(environments.Links.Next.Href),
			Description: types.StringValue(environments.Links.Next.Description),
			BasePath:    types.StringValue(environments.Links.Next.BasePath),
		},
		Previous: models.EnvironmentsLinkModel{
			Href:        types.StringValue(environments.Links.Previous.Href),
			Description: types.StringValue(environments.Links.Previous.Description),
			BasePath:    types.StringValue(environments.Links.Previous.BasePath),
		},
		Last: models.EnvironmentsLinkModel{
			Href:        types.StringValue(environments.Links.Last.Href),
			Description: types.StringValue(environments.Links.Last.Description),
			BasePath:    types.StringValue(environments.Links.Last.BasePath),
		},
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

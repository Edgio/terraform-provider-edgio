package resources

import (
	"context"
	"fmt"
	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var _ resource.Resource = &PropertyResource{}

type PropertyResource struct {
	client *edgio_api.EdgioClient
}

func NewPropertyResource(client *edgio_api.EdgioClient) resource.Resource {
	return &PropertyResource{
		client: client,
	}
}

func (r *PropertyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property"
}

func (r *PropertyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"slug": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *PropertyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PropertyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	property, err := r.client.CreateProperty(ctx, plan.OrganizationID.ValueString(), plan.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Property",
			fmt.Sprintf("Could not create property, unexpected error: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(property.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PropertyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PropertyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	property, err := r.client.GetSpecificProperty(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Property",
			fmt.Sprintf("Could not read property ID %s, unexpected error: %s", state.ID.ValueString(), err),
		)
		return
	}
	state.OrganizationID = types.StringValue(property.OrganizationID)
	state.Slug = types.StringValue(property.Slug)
	state.Type = types.StringValue(property.Type)
	state.CreatedAt = types.StringValue(property.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(property.UpdatedAt.Format(time.RFC3339))

	state.OrganizationID = types.StringValue(property.OrganizationID)
	state.Slug = types.StringValue(property.Slug)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *PropertyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// If your API doesn't support updates, you can return an error here
	resp.Diagnostics.AddError(
		"Error Updating Property",
		"Property update is not supported",
	)
}

func (r *PropertyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PropertyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to delete the property
	// If your API doesn't support deletion, you might want to handle this differently
	err := r.client.DeleteProperty(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Property",
			fmt.Sprintf("Could not delete property, unexpected error: %s", err),
		)
		return
	}
}

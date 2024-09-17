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
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"slug": schema.StringAttribute{
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"id_link": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
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

	plan.OrganizationID = types.StringValue(property.OrganizationID)
	plan.Slug = types.StringValue(property.Slug)
	plan.Id = types.StringValue(property.Id)
	plan.IdLink = types.StringValue(property.IdLink)
	plan.Type = types.StringValue(property.Type)
	plan.CreatedAt = types.StringValue(property.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(property.UpdatedAt.Format(time.RFC3339))

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

	property, err := r.client.GetProperty(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Property",
			fmt.Sprintf("Could not read property ID %s, unexpected error: %s", state.Id.ValueString(), err),
		)
		return
	}
	state.OrganizationID = types.StringValue(property.OrganizationID)
	state.Slug = types.StringValue(property.Slug)
	state.Id = types.StringValue(property.Id)
	state.IdLink = types.StringValue(property.IdLink)
	state.Type = types.StringValue(property.Type)
	state.CreatedAt = types.StringValue(property.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(property.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *PropertyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.PropertyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan models.PropertyModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProperty, err := r.client.UpdateProperty(ctx, state.Id.ValueString(), plan.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating property",
			fmt.Sprintf("Could not update property ID %s, unexpected error: %s", state.Id.String(), err),
		)
		return
	}

	plan.OrganizationID = types.StringValue(updatedProperty.OrganizationID)
	plan.Slug = types.StringValue(updatedProperty.Slug)
	plan.Id = types.StringValue(updatedProperty.Id)
	plan.IdLink = types.StringValue(updatedProperty.IdLink)
	plan.Type = types.StringValue(updatedProperty.Type)
	plan.CreatedAt = types.StringValue(updatedProperty.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(updatedProperty.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PropertyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PropertyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProperty(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Property",
			fmt.Sprintf("Could not delete property, unexpected error: %s", err),
		)
		return
	}
}

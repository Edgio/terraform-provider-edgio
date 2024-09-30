package resources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_api/dtos/purge"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the resource.Resource interface
var _ resource.Resource = &PurgeCacheResource{}

type PurgeCacheResource struct {
	client *edgio_api.EdgioClient
}

func NewPurgeCacheResource(client *edgio_api.EdgioClient) resource.Resource {
	return &PurgeCacheResource{
		client: client,
	}
}

func (r *PurgeCacheResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "edgio_purge_cache"
}

func (r *PurgeCacheResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to trigger a cache purge for a given environment.",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Description: "ID of the environment for which to purge cache.",
				Required:    true,
			},
			"purge_type": schema.StringAttribute{
				Description: `Defines the set of content to be purged. Available values are:
				- all_entries: Purge all cached content.
				- path: Purge one or more relative path(s).
				- surrogate_key: Purge one or more surrogate key(s).`,
				Required: true,
			},
			"values": schema.ListAttribute{
				Description: `Defines the content to be purged based on the purge_type.
				- all_entries: Omit this property or pass an empty array.
				- path: Provide one or more relative path(s).
				- surrogate_key: Provide one or more surrogate key(s).`,
				ElementType: types.StringType,
				Required:    true,
			},
			"hostname": schema.StringAttribute{
				Description: `Optional. If specified, cached paths will be purged for this specific hostname. If omitted, the specified paths will be purged for all hostnames.`,
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique ID of the purge operation.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the purge operation.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp of when the purge request was created.",
				Computed:    true,
			},
			"completed_at": schema.StringAttribute{
				Description: "Timestamp of when the purge request was completed.",
				Computed:    true,
			},
			"progress_percentage": schema.Float32Attribute{
				Description: "Percentage progress of the purge operation.",
				Computed:    true,
			},
		},
	}
}

func (r *PurgeCacheResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PurgeCacheResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	purgeResponse, err := r.client.GetPurgeStatus(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Purge Status Error", err.Error())
		return
	}

	state.ID = types.StringValue(purgeResponse.ID)
	state.Status = types.StringValue(purgeResponse.Status)
	state.CreatedAt = types.StringValue(purgeResponse.CreatedAt.Format(time.RFC3339))
	state.CompletedAt = types.StringValue(purgeResponse.CompletedAt.Format(time.RFC3339))
	state.ProgressPercentage = types.Float32Value(purgeResponse.ProgressPercentage)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *PurgeCacheResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PurgeCacheResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var values []string
	diags = plan.Values.ElementsAs(ctx, &values, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	purgeRequest := purge.PurgeRequest{
		EnvironmentID: plan.EnvironmentID.ValueString(),
		PurgeType:     plan.PurgeType.ValueString(),
		Values:        values,
		Hostname:      plan.Hostname.ValueStringPointer(),
	}

	purgeResponse, err := r.client.PurgeCache(&purgeRequest)
	if err != nil {
		resp.Diagnostics.AddError("Purge Cache Error", err.Error())
		return
	}

	state := models.PurgeCacheResourceModel{
		ID:                 types.StringValue(purgeResponse.ID),
		Status:             types.StringValue(purgeResponse.Status),
		CreatedAt:          types.StringValue(purgeResponse.CreatedAt.Format(time.RFC3339)),
		CompletedAt:        types.StringValue(purgeResponse.CompletedAt.Format(time.RFC3339)),
		ProgressPercentage: types.Float32Value(purgeResponse.ProgressPercentage),
		Values:             plan.Values,
		EnvironmentID:      plan.EnvironmentID,
		PurgeType:          plan.PurgeType,
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PurgeCacheResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Error Updating Purge Cache", "Purge Cache cannot be updated")
}

func (r *PurgeCacheResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Error Deleting Purge Cache", "Purge Cache cannot be deleted")
}

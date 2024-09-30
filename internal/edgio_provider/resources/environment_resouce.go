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
var _ resource.Resource = &EnvironmentResource{}

type EnvironmentResource struct {
	client *edgio_api.EdgioClient
}

func NewEnvironmentResource(client *edgio_api.EdgioClient) *EnvironmentResource {
	return &EnvironmentResource{
		client: client,
	}
}

func (r *EnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "edgio_environment"
}

func (r *EnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"property_id": schema.StringAttribute{
				Required: true,
			},
			"legacy_account_number": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"can_members_deploy": schema.BoolAttribute{
				Optional: true,
			},
			"only_maintainers_can_deploy": schema.BoolAttribute{
				Optional: true,
			},
			"http_request_logging": schema.BoolAttribute{
				Optional: true,
			},
			"default_domain_name": schema.StringAttribute{
				Computed: true,
			},
			"pci_compliance": schema.BoolAttribute{
				Computed: true,
			},
			"dns_domain_name": schema.StringAttribute{
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

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.EnvironmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.CreateEnvironment(plan.PropertyID.ValueString(), plan.Name.ValueString(), plan.CanMembersDeploy.ValueBool(), plan.OnlyMaintainersCanDeploy.ValueBool(), plan.HttpRequestLogging.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Environment",
			fmt.Sprintf("Could not create environment, unexpected error: %s", err),
		)
		return
	}

	plan.Id = types.StringValue(env.Id)
	plan.PropertyID = types.StringValue(env.PropertyID)
	plan.LegacyAccountNumber = types.StringValue(env.LegacyAccountNumber)
	plan.Name = types.StringValue(env.Name)
	plan.CanMembersDeploy = types.BoolValue(env.CanMembersDeploy)
	plan.OnlyMaintainersCanDeploy = types.BoolValue(env.OnlyMaintainersCanDeploy)
	plan.HttpRequestLogging = types.BoolValue(env.HttpRequestLogging)
	plan.DefaultDomainName = types.StringValue(env.DefaultDomainName)
	plan.PciCompliance = types.BoolValue(env.PciCompliance)
	plan.DnsDomainName = types.StringValue(env.DnsDomainName)
	plan.CreatedAt = types.StringValue(env.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(env.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.EnvironmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.GetEnvironment(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Environment",
			fmt.Sprintf("Could not read environment ID %s, unexpected error: %s", state.Id.ValueString(), err),
		)
		return
	}

	state.PropertyID = types.StringValue(env.PropertyID)
	state.LegacyAccountNumber = types.StringValue(env.LegacyAccountNumber)
	state.Name = types.StringValue(env.Name)
	state.CanMembersDeploy = types.BoolValue(env.CanMembersDeploy)
	state.OnlyMaintainersCanDeploy = types.BoolValue(env.OnlyMaintainersCanDeploy)
	state.HttpRequestLogging = types.BoolValue(env.HttpRequestLogging)
	state.DefaultDomainName = types.StringValue(env.DefaultDomainName)
	state.PciCompliance = types.BoolValue(env.PciCompliance)
	state.DnsDomainName = types.StringValue(env.DnsDomainName)
	state.CreatedAt = types.StringValue(env.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(env.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.EnvironmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan models.EnvironmentModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedEnv, err := r.client.UpdateEnvironment(
		state.Id.ValueString(),
		plan.Name.ValueString(),
		plan.CanMembersDeploy.ValueBool(),
		plan.HttpRequestLogging.ValueBool(),
		// TODO: What to do with this, it's not in the plan?
		false)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment",
			fmt.Sprintf("Could not update environment ID %s, unexpected error: %s", state.Id.ValueString(), err),
		)
		return
	}

	plan.Id = types.StringValue(updatedEnv.Id)
	plan.PropertyID = types.StringValue(updatedEnv.PropertyID)
	plan.LegacyAccountNumber = types.StringValue(updatedEnv.LegacyAccountNumber)
	plan.Name = types.StringValue(updatedEnv.Name)
	plan.CanMembersDeploy = types.BoolValue(updatedEnv.CanMembersDeploy)
	plan.OnlyMaintainersCanDeploy = types.BoolValue(updatedEnv.OnlyMaintainersCanDeploy)
	plan.HttpRequestLogging = types.BoolValue(updatedEnv.HttpRequestLogging)
	plan.DefaultDomainName = types.StringValue(updatedEnv.DefaultDomainName)
	plan.PciCompliance = types.BoolValue(updatedEnv.PciCompliance)
	plan.DnsDomainName = types.StringValue(updatedEnv.DnsDomainName)
	plan.CreatedAt = types.StringValue(updatedEnv.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(updatedEnv.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.EnvironmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEnvironment(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Environment",
			fmt.Sprintf("Could not delete environment, unexpected error: %s", err),
		)
		return
	}
}

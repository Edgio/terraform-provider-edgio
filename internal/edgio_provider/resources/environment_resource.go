package resources

import (
	"context"
	"fmt"
	"terraform-provider-edgio/internal/edgio_api"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"terraform-provider-edgio/internal/edgio_provider/utility"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &EnvironmentResource{}

type EnvironmentResource struct {
	client edgio_api.EdgioClientInterface
}

func NewEnvironmentResource(client edgio_api.EdgioClientInterface) *EnvironmentResource {
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
			"only_maintainers_can_deploy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"http_request_logging": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"default_domain_name": schema.StringAttribute{
				Computed: true,
			},
			"pci_compliance": schema.BoolAttribute{
				Computed: true,
				Default:  booldefault.StaticBool(false),
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

	env, err := r.client.CreateEnvironment(
		plan.PropertyID.ValueString(),
		plan.Name.ValueString(),
		plan.OnlyMaintainersCanDeploy.ValueBool(),
		plan.HttpRequestLogging.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Environment",
			fmt.Sprintf("Could not create environment, unexpected error: %s", err),
		)
		return
	}

	newState := utility.ConvertEnvironmentToModel(env)
	diags = resp.State.Set(ctx, newState)
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

	newState := utility.ConvertEnvironmentToModel(env)
	diags = resp.State.Set(ctx, newState)
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
		plan.OnlyMaintainersCanDeploy.ValueBool(),
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

	newState := utility.ConvertEnvironmentToModel(updatedEnv)
	diags = resp.State.Set(ctx, newState)
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

package resources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_provider/utility"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// JSONEqualityModifier is a custom plan modifier that formats JSON strings always the same
type JSONEqualityModifier struct{}

// Description provides a description of the modifier
func (m JSONEqualityModifier) Description(ctx context.Context) string {
	return "Formats JSON strings always the same, ignoring formatting differences"
}

// MarkdownDescription provides a description of the modifier in markdown format
func (m JSONEqualityModifier) MarkdownDescription(ctx context.Context) string {
	return "Formats JSON strings always the same, ignoring formatting differences"
}

// PlanModifyString is called to modify the planned new state
func (m JSONEqualityModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the plan is unknown, we can't make any decisions
	if req.PlanValue.IsUnknown() {
		return
	}

	json, error := utility.MinifyJSON(req.PlanValue.ValueString())

	if error != nil {
		resp.Diagnostics.AddError("Error parsing plan JSON", error.Error())
		return
	}

	resp.PlanValue = types.StringValue(json)
}

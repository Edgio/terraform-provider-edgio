package resources

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NormalizeJsonPlanModifier struct{}

func (m NormalizeJsonPlanModifier) Description(ctx context.Context) string {
	return "Normalize JSON values in the plan."
}

func (m NormalizeJsonPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Normalize JSON values in the plan."
}

func (m NormalizeJsonPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	planValue, diag := req.PlanValue.ToStringValue(ctx)
	if diag != nil {
		resp.Diagnostics.Append(diag...)
		return
	}

	normalized, err := normalizeJSON(planValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to normalize JSON value", err.Error())
		return
	}

	resp.PlanValue = types.StringValue(normalized)
}

func normalizeJSON(input string) (string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return "", err
	}

	normalizedBytes, err := json.Marshal(raw)
	if err != nil {
		return "", err
	}

	return string(normalizedBytes), nil
}

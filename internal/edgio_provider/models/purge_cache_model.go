package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PurgeCacheResourceModel struct {
	EnvironmentID      types.String  `tfsdk:"environment_id"`
	PurgeType          types.String  `tfsdk:"purge_type"`
	Values             types.List    `tfsdk:"values"`
	Hostname           types.String  `tfsdk:"hostname"`
	ID                 types.String  `tfsdk:"id"`
	Status             types.String  `tfsdk:"status"`
	CreatedAt          types.String  `tfsdk:"created_at"`
	CompletedAt        types.String  `tfsdk:"completed_at"`
	ProgressPercentage types.Float32 `tfsdk:"progress_percentage"`
}

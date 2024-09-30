package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnvironmentsModel struct {
	Type         types.String           `tfsdk:"type"`
	Id           types.String           `tfsdk:"id"`
	Links        EnvironmentsLinksModel `tfsdk:"links"`
	TotalItems   types.Int32            `tfsdk:"total_items"`
	Environments []EnvironmentModel     `tfsdk:"items"`
}

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnvironmentsLinkModel struct {
	Href        types.String `tfsdk:"href"`
	Description types.String `tfsdk:"description"`
	BasePath    types.String `tfsdk:"base_path"`
}

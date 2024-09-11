package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PropertyModel struct {
	Type           types.String `tfsdk:"type"`
	ID             types.String `tfsdk:"id"`
	PropertyID     types.String `tfsdk:"property_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Slug           types.String `tfsdk:"slug"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

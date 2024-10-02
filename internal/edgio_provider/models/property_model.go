package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PropertyModel struct {
	Type           types.String `tfsdk:"type"`
	Id             types.String `tfsdk:"id"`
	IdLink         types.String `tfsdk:"id_link"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Slug           types.String `tfsdk:"slug"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

type PropertiesModel struct {
	OrganizationID types.String    `tfsdk:"organization_id"`
	ItemCount      types.Int32     `tfsdk:"item_count"`
	Properties     []PropertyModel `tfsdk:"properties"`
}

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PropertiesModel struct {
	OrganizationID types.String         `tfsdk:"organization_id"`
	ItemCount      types.Int32          `tfsdk:"item_count"`
	Properties     []PropertyModel      `tfsdk:"properties"`
	Links          PropertiesLinksModel `tfsdk:"links"`
}

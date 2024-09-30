package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TLSCertsModel struct {
	EnvironmentID types.String   `tfsdk:"environment_id"`
	Page          types.Int32    `tfsdk:"page"`
	PageSize      types.Int32    `tfsdk:"page_size"`
	ItemCount     types.Int32    `tfsdk:"item_count"`
	Certificates  []TLSCertModel `tfsdk:"certificates"`
}

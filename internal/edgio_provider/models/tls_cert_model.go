package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TLSCertModel struct {
	EnvironmentID    types.String `tfsdk:"environment_id"`
	PrimaryCert      types.String `tfsdk:"primary_cert"`
	IntermediateCert types.String `tfsdk:"intermediate_cert"`
	PrivateKey       types.String `tfsdk:"private_key"`
	ID               types.String `tfsdk:"id"`
	Expiration       types.String `tfsdk:"expiration"`
	Status           types.String `tfsdk:"status"`
	Generated        types.Bool   `tfsdk:"generated"`
	Serial           types.String `tfsdk:"serial"`
	CommonName       types.String `tfsdk:"common_name"`
	AlternativeNames types.List   `tfsdk:"alternative_names"`
	ActivationError  types.String `tfsdk:"activation_error"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

type TLSCertsModel struct {
	EnvironmentID types.String   `tfsdk:"environment_id"`
	Page          types.Int32    `tfsdk:"page"`
	PageSize      types.Int32    `tfsdk:"page_size"`
	ItemCount     types.Int32    `tfsdk:"item_count"`
	Certificates  []TLSCertModel `tfsdk:"certificates"`
}

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvironmentModel struct {
	Id                  types.String `tfsdk:"id"`
	PropertyID          types.String `tfsdk:"property_id"`
	LegacyAccountNumber types.String `tfsdk:"legacy_account_number"`
	Name                types.String `tfsdk:"name"`
	CanMembersDeploy    types.Bool   `tfsdk:"can_members_deploy"`
	// todo: uncomment this when the API is updated
	// OnlyMaintainersCanDeploy types.Bool   `tfsdk:"only_maintainers_can_deploy"`
	HttpRequestLogging types.Bool   `tfsdk:"http_request_logging"`
	DefaultDomainName  types.String `tfsdk:"default_domain_name"`
	PciCompliance      types.Bool   `tfsdk:"pci_compliance"`
	DnsDomainName      types.String `tfsdk:"dns_domain_name"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

type EnvironmentsModel struct {
	Type         types.String       `tfsdk:"type"`
	Id           types.String       `tfsdk:"id"`
	TotalItems   types.Int32        `tfsdk:"total_items"`
	Environments []EnvironmentModel `tfsdk:"items"`
}

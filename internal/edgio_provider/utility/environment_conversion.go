package utility

import (
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertEnvironmentToModel(env *dtos.Environment) models.EnvironmentModel {
	return models.EnvironmentModel{
		Id:                  types.StringValue(env.Id),
		PropertyID:          types.StringValue(env.PropertyID),
		LegacyAccountNumber: types.StringValue(env.LegacyAccountNumber),
		Name:                types.StringValue(env.Name),
		CanMembersDeploy:    types.BoolValue(env.CanMembersDeploy),
		// OnlyMaintainersCanDeploy: types.BoolValue(env.OnlyMaintainersCanDeploy),
		HttpRequestLogging: types.BoolValue(env.HttpRequestLogging),
		DefaultDomainName:  types.StringValue(env.DefaultDomainName),
		PciCompliance:      types.BoolValue(env.PciCompliance),
		DnsDomainName:      types.StringValue(env.DnsDomainName),
		CreatedAt:          types.StringValue(env.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:          types.StringValue(env.UpdatedAt.Format(time.RFC3339)),
	}
}

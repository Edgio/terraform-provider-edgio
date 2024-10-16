package utility

import (
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertPropertyToModel(property *dtos.Property) models.PropertyModel {
	return models.PropertyModel{
		OrganizationID: types.StringValue(property.OrganizationID),
		Slug:           types.StringValue(property.Slug),
		Id:             types.StringValue(property.Id),
		IdLink:         types.StringValue(property.IdLink),
		CreatedAt:      types.StringValue(property.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:      types.StringValue(property.UpdatedAt.Format(time.RFC3339)),
	}
}

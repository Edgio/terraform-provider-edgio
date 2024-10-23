package utility

import (
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertTlsCertsToModel(tlsRes *dtos.TLSCertResponse) models.TLSCertModel {
	return models.TLSCertModel{
		ID:               types.StringValue(tlsRes.ID),
		EnvironmentID:    types.StringValue(tlsRes.EnvironmentID),
		PrimaryCert:      types.StringValue(tlsRes.PrimaryCert),
		IntermediateCert: types.StringValue(tlsRes.IntermediateCert),
		Expiration:       types.StringValue(tlsRes.Expiration),
		Status:           types.StringValue(tlsRes.Status),
		Generated:        types.BoolValue(tlsRes.Generated),
		Serial:           types.StringValue(tlsRes.Serial),
		CommonName:       types.StringValue(tlsRes.CommonName),
		AlternativeNames: StringSliceToTypesList(&tlsRes.AlternativeNames),
		ActivationError:  types.StringValue(tlsRes.ActivationError),
		CreatedAt:        types.StringValue(tlsRes.CreatedAt),
		UpdatedAt:        types.StringValue(tlsRes.UpdatedAt),
	}
}

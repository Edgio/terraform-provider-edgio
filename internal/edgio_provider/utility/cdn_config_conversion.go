package utility

import (
	"encoding/json"
	"terraform-provider-edgio/internal/edgio_api/dtos"
	"terraform-provider-edgio/internal/edgio_provider/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertCdnConfigToNative(model *models.CDNConfigurationModel) dtos.CDNConfiguration {
	rules := json.RawMessage(model.Rules.ValueString())

	return dtos.CDNConfiguration{
		ConfigurationID:        model.ConfigurationID.ValueString(),
		Rules:                  rules,
		EnvironmentID:          model.EnvironmentID.ValueString(),
		Origins:                convertOriginsToNative(model.Origins),
		Hostnames:              convertHostnamesToNative(model.Hostnames),
		Experiments:            TypesListToStringSlicePointer(model.Experiments),
		EdgeFunctionsSources:   MapValueToStringMapPointer(model.EdgeFunctionsSources),
		EdgeFunctionInitScript: ToPtrString(model.EdgeFunctionInitScript),
	}
}

func ConvertNativeToCdnConfig(dto *dtos.CDNConfiguration) models.CDNConfigurationModel {
	rulesStr := ""

	if dto.Rules != nil {
		// failing to parse will lead to empty rulesStr, which
		// will fail on an endpoint, so that's why we ignore errors here
		rulesBytes, _ := dto.Rules.MarshalJSON()
		rulesStr = string(rulesBytes)
		str, _ := MinifyJSON(rulesStr)
		rulesStr = str
	}

	return models.CDNConfigurationModel{
		ConfigurationID:        types.StringValue(dto.ConfigurationID),
		EnvironmentID:          types.StringValue(dto.EnvironmentID),
		Origins:                convertNativeToOrigins(dto.Origins),
		Rules:                  types.StringValue(rulesStr),
		Hostnames:              convertNativeToHostnames(dto.Hostnames),
		Experiments:            StringSliceToTypesList(dto.Experiments),
		EdgeFunctionsSources:   StringMapToMapValue(dto.EdgeFunctionsSources),
		EdgeFunctionInitScript: types.StringPointerValue(dto.EdgeFunctionInitScript),
	}
}

func convertOriginsToNative(origins []models.OriginModel) []dtos.Origin {
	var natives []dtos.Origin
	for _, origin := range origins {
		natives = append(natives, dtos.Origin{
			Name:                origin.Name.ValueString(),
			Type:                ToPtrString(origin.Type),
			Hosts:               convertHostsToNative(origin.Hosts),
			Balancer:            ToPtrString(origin.Balancer),
			OverrideHostHeader:  ToPtrString(origin.OverrideHostHeader),
			Shields:             convertShieldsToNative(origin.Shields),
			PciCertifiedShields: ToPtrBool(origin.PciCertifiedShields),
			TLSVerify:           convertTLSVerifyToNative(origin.TLSVerify),
			Retry:               convertRetryToNative(origin.Retry),
		})
	}
	return natives
}

func convertNativeToOrigins(origins []dtos.Origin) []models.OriginModel {
	var m []models.OriginModel
	for _, origin := range origins {
		m = append(m, models.OriginModel{
			Name:                types.StringValue(origin.Name),
			Type:                types.StringPointerValue(origin.Type),
			Hosts:               convertNativeToHosts(origin.Hosts),
			Balancer:            types.StringPointerValue(origin.Balancer),
			OverrideHostHeader:  types.StringPointerValue(origin.OverrideHostHeader),
			Shields:             convertNativeToShields(origin.Shields),
			PciCertifiedShields: types.BoolPointerValue(origin.PciCertifiedShields),
			TLSVerify:           convertNativeToTLSVerify(origin.TLSVerify),
			Retry:               convertNativeToRetry(origin.Retry),
		})
	}
	return m
}

func convertHostnamesToNative(hostnames []models.HostnameModel) []dtos.Hostname {
	var natives []dtos.Hostname
	for _, hostname := range hostnames {
		natives = append(natives, dtos.Hostname{
			Hostname:          ToPtrString(hostname.Hostname),
			DefaultOriginName: ToPtrString(hostname.DefaultOriginName),
			ReportCode:        ToPtrInt64(hostname.ReportCode),
			TLS:               convertTLSToNative(hostname.TLS),
			Directory:         ToPtrString(hostname.Directory),
		})
	}
	return natives
}

func convertNativeToHostnames(hostnames []dtos.Hostname) []models.HostnameModel {
	var natives []models.HostnameModel
	for _, hostname := range hostnames {
		natives = append(natives, models.HostnameModel{
			Hostname:          types.StringPointerValue(hostname.Hostname),
			DefaultOriginName: types.StringPointerValue(hostname.DefaultOriginName),
			ReportCode:        types.Int64PointerValue(hostname.ReportCode),
			TLS:               convertNativeToTLS(hostname.TLS),
			Directory:         types.StringPointerValue(hostname.Directory),
		})
	}
	return natives
}

func convertTLSToNative(tls *models.TLSModel) *dtos.TLS {
	if tls == nil {
		return nil
	}

	return &dtos.TLS{
		NPN:                 ToPtrBool(tls.NPN),
		ALPN:                ToPtrBool(tls.ALPN),
		Protocols:           ToPtrString(tls.Protocols),
		UseSigAlgs:          ToPtrBool(tls.UseSigAlgs),
		SNI:                 ToPtrBool(tls.SNI),
		SniStrict:           ToPtrBool(tls.SniStrict),
		SniHostMatch:        ToPtrBool(tls.SniHostMatch),
		ClientRenegotiation: ToPtrBool(tls.ClientRenegotiation),
		Options:             ToPtrString(tls.Options),
		CipherList:          ToPtrString(tls.CipherList),
		NamedCurve:          ToPtrString(tls.NamedCurve),
		OCSP:                ToPtrBool(tls.OCSP),
		PEM:                 ToPtrString(tls.PEM),
		CA:                  ToPtrString(tls.CA),
	}
}

func convertNativeToTLS(tls *dtos.TLS) *models.TLSModel {
	if tls == nil {
		return nil
	}

	return &models.TLSModel{
		NPN:                 types.BoolPointerValue(tls.NPN),
		ALPN:                types.BoolPointerValue(tls.ALPN),
		Protocols:           types.StringPointerValue(tls.Protocols),
		UseSigAlgs:          types.BoolPointerValue(tls.UseSigAlgs),
		SNI:                 types.BoolPointerValue(tls.SNI),
		SniStrict:           types.BoolPointerValue(tls.SniStrict),
		SniHostMatch:        types.BoolPointerValue(tls.SniHostMatch),
		ClientRenegotiation: types.BoolPointerValue(tls.ClientRenegotiation),
		Options:             types.StringPointerValue(tls.Options),
		CipherList:          types.StringPointerValue(tls.CipherList),
		NamedCurve:          types.StringPointerValue(tls.NamedCurve),
		OCSP:                types.BoolPointerValue(tls.OCSP),
		PEM:                 types.StringPointerValue(tls.PEM),
		CA:                  types.StringPointerValue(tls.CA),
	}
}

func convertRetryToNative(retry *models.RetryModel) *dtos.Retry {
	if retry == nil {
		return nil
	}

	return &dtos.Retry{
		StatusCodes:            TypesListToIntSlicePointer(retry.StatusCodes),
		IgnoreRetryAfterHeader: ToPtrBool(retry.IgnoreRetryAfterHeader),
		AfterSeconds:           ToPtrInt64(retry.AfterSeconds),
		MaxRequests:            ToPtrInt64(retry.MaxRequests),
		MaxWaitSeconds:         ToPtrInt64(retry.MaxWaitSeconds),
	}
}

func convertNativeToRetry(retry *dtos.Retry) *models.RetryModel {
	if retry == nil {
		return nil
	}

	return &models.RetryModel{
		StatusCodes:            IntSliceToTypesList(retry.StatusCodes),
		IgnoreRetryAfterHeader: types.BoolPointerValue(retry.IgnoreRetryAfterHeader),
		AfterSeconds:           types.Int64PointerValue(retry.AfterSeconds),
		MaxRequests:            types.Int64PointerValue(retry.MaxRequests),
		MaxWaitSeconds:         types.Int64PointerValue(retry.MaxWaitSeconds),
	}
}

func convertShieldsToNative(shields *models.ShieldsModel) *dtos.Shields {
	if shields == nil {
		return nil
	}

	return &dtos.Shields{
		Apac:   ToPtrString(shields.Apac),
		Emea:   ToPtrString(shields.Emea),
		USWest: ToPtrString(shields.USWest),
		USEast: ToPtrString(shields.USEast),
	}
}

func convertNativeToShields(shields *dtos.Shields) *models.ShieldsModel {
	if shields == nil {
		return nil
	}

	return &models.ShieldsModel{
		Apac:   types.StringPointerValue(shields.Apac),
		Emea:   types.StringPointerValue(shields.Emea),
		USWest: types.StringPointerValue(shields.USWest),
		USEast: types.StringPointerValue(shields.USEast),
	}
}

func convertTLSVerifyToNative(tlsVerify *models.TLSVerifyModel) *dtos.TLSVerify {
	if tlsVerify == nil {
		return nil
	}

	return &dtos.TLSVerify{
		UseSNI:                   ToPtrBool(tlsVerify.UseSNI),
		SNIHintAndStrictSanCheck: ToPtrString(tlsVerify.SNIHintAndStrictSanCheck),
		AllowSelfSignedCerts:     ToPtrBool(tlsVerify.AllowSelfSignedCerts),
		PinnedCerts:              TypesListToStringSlicePointer(tlsVerify.PinnedCerts),
	}
}

func convertNativeToTLSVerify(tlsVerify *dtos.TLSVerify) *models.TLSVerifyModel {
	if tlsVerify == nil {
		return nil
	}

	return &models.TLSVerifyModel{
		UseSNI:                   types.BoolPointerValue(tlsVerify.UseSNI),
		SNIHintAndStrictSanCheck: types.StringPointerValue(tlsVerify.SNIHintAndStrictSanCheck),
		AllowSelfSignedCerts:     types.BoolPointerValue(tlsVerify.AllowSelfSignedCerts),
		PinnedCerts:              StringSliceToTypesList(tlsVerify.PinnedCerts),
	}
}

func convertHostsToNative(hosts []models.HostModel) []dtos.Host {
	var natives []dtos.Host
	for _, host := range hosts {
		natives = append(natives, dtos.Host{
			Weight:                   ToPtrInt64(host.Weight),
			DNSMaxTTL:                ToPtrInt64(host.DNSMaxTTL),
			DNSPreference:            ToPtrString(host.DNSPreference),
			MaxHardPool:              ToPtrInt64(host.MaxHardPool),
			DNSMinTTL:                ToPtrInt64(host.DNSMinTTL),
			Location:                 convertLocationToNative(host.Location),
			MaxPool:                  ToPtrInt64(host.MaxPool),
			Balancer:                 ToPtrString(host.Balancer),
			Scheme:                   ToPtrString(host.Scheme),
			OverrideHostHeader:       ToPtrString(host.OverrideHostHeader),
			SNIHintAndStrictSanCheck: ToPtrString(host.SNIHintAndStrictSanCheck),
			UseSNI:                   ToPtrBool(host.UseSNI),
		})
	}
	return natives
}

func convertNativeToHosts(hosts []dtos.Host) []models.HostModel {
	var m []models.HostModel
	for _, host := range hosts {
		m = append(m, models.HostModel{
			Weight:                   types.Int64PointerValue(host.Weight),
			DNSMaxTTL:                types.Int64PointerValue(host.DNSMaxTTL),
			DNSPreference:            types.StringPointerValue(host.DNSPreference),
			MaxHardPool:              types.Int64PointerValue(host.MaxHardPool),
			DNSMinTTL:                types.Int64PointerValue(host.DNSMinTTL),
			Location:                 convertNativeToLocation(host.Location),
			MaxPool:                  types.Int64PointerValue(host.MaxPool),
			Balancer:                 types.StringPointerValue(host.Balancer),
			Scheme:                   types.StringPointerValue(host.Scheme),
			OverrideHostHeader:       types.StringPointerValue(host.OverrideHostHeader),
			SNIHintAndStrictSanCheck: types.StringPointerValue(host.SNIHintAndStrictSanCheck),
			UseSNI:                   types.BoolPointerValue(host.UseSNI),
		})
	}
	return m
}

func convertLocationToNative(locations *[]models.LocationModel) *[]dtos.Location {
	if locations == nil {
		return nil
	}

	var natives []dtos.Location
	for _, location := range *locations {
		natives = append(natives, dtos.Location{
			Port:     ToPtrInt64(location.Port),
			Hostname: ToPtrString(location.Hostname),
		})
	}
	return &natives
}

func convertNativeToLocation(locations *[]dtos.Location) *[]models.LocationModel {
	if locations == nil {
		return nil
	}

	var m []models.LocationModel
	for _, location := range *locations {
		m = append(m, models.LocationModel{
			Port:     types.Int64PointerValue(location.Port),
			Hostname: types.StringPointerValue(location.Hostname),
		})
	}
	return &m
}

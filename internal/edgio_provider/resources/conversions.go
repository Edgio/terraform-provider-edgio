package resources

import (
	"context"
	"encoding/json"
	"terraform-provider-edgio/internal/edgio_api/dtos/configuration"
	"terraform-provider-edgio/internal/edgio_provider/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertCdnConfigToNative(model *models.CDNConfigurationModel) configuration.CDNConfiguration {
	rules := json.RawMessage(model.Rules.ValueString())

	return configuration.CDNConfiguration{
		ConfigurationID:        model.ConfigurationID.ValueString(),
		Rules:                  rules,
		EnvironmentID:          model.EnvironmentID.ValueString(),
		Origins:                convertOriginsToNative(model.Origins),
		Hostnames:              convertHostnamesToNative(model.Hostnames),
		Experiments:            TypesListToStringSlice(model.Experiments),
		EdgeFunctionsSources:   MapValueToStringMap(model.EdgeFunctionsSources),
		EdgeFunctionInitScript: model.EdgeFunctionInitScript.ValueString(),
	}
}

func ConvertNativeToCdnConfig(dto *configuration.CDNConfiguration) models.CDNConfigurationModel {
	rulesStr := ""
	if dto.Rules != nil {
		// failing to parse will lead to empty rulesStr, which
		// will fail on an endpoint, so that's why we ignore errors here
		rulesBytes, _ := json.Marshal(dto.Rules)
		rulesStr = string(rulesBytes)
	}

	return models.CDNConfigurationModel{
		ConfigurationID:        types.StringValue(dto.ConfigurationID),
		EnvironmentID:          types.StringValue(dto.EnvironmentID),
		Origins:                convertNativeToOrigins(dto.Origins),
		Rules:                  types.StringValue(rulesStr),
		Hostnames:              convertNativeToHostnames(dto.Hostnames),
		Experiments:            StringSliceToTypesList(dto.Experiments),
		EdgeFunctionsSources:   StringMapToMapValue(dto.EdgeFunctionsSources),
		EdgeFunctionInitScript: types.StringValue(dto.EdgeFunctionInitScript),
	}
}

func convertOriginsToNative(origins []models.OriginModel) []configuration.Origin {
	var natives []configuration.Origin
	for _, origin := range origins {
		natives = append(natives, configuration.Origin{
			Name:                origin.Name.ValueString(),
			Type:                origin.Type.ValueString(),
			Hosts:               convertHostsToNative(origin.Hosts),
			Balancer:            origin.Balancer.ValueString(),
			OverrideHostHeader:  origin.OverrideHostHeader.ValueString(),
			Shields:             convertShieldsToNative(origin.Shields),
			PciCertifiedShields: origin.PciCertifiedShields.ValueBool(),
			TLSVerify:           convertTLSVerifyToNative(origin.TLSVerify),
			Retry:               convertRetryToNative(origin.Retry),
		})
	}
	return natives
}

func convertNativeToOrigins(origins []configuration.Origin) []models.OriginModel {
	var m []models.OriginModel
	for _, origin := range origins {
		m = append(m, models.OriginModel{
			Name:                types.StringValue(origin.Name),
			Type:                types.StringValue(origin.Type),
			Hosts:               convertNativeToHosts(origin.Hosts),
			Balancer:            types.StringValue(origin.Balancer),
			OverrideHostHeader:  types.StringValue(origin.OverrideHostHeader),
			Shields:             convertNativeToShields(origin.Shields),
			PciCertifiedShields: types.BoolValue(origin.PciCertifiedShields),
			TLSVerify:           convertNativeToTLSVerify(origin.TLSVerify),
			Retry:               convertNativeToRetry(origin.Retry),
		})
	}
	return m
}

func convertHostnamesToNative(hostnames []models.HostnameModel) []configuration.Hostname {
	var natives []configuration.Hostname
	for _, hostname := range hostnames {
		natives = append(natives, configuration.Hostname{
			Hostname:          hostname.Hostname.ValueString(),
			DefaultOriginName: hostname.DefaultOriginName.ValueString(),
			ReportCode:        int(hostname.ReportCode.ValueInt64()),
			TLS:               convertTLSToNative(hostname.TLS),
			Directory:         hostname.Directory.ValueString(),
		})
	}
	return natives
}

func convertNativeToHostnames(hostnames []configuration.Hostname) []models.HostnameModel {
	var natives []models.HostnameModel
	for _, hostname := range hostnames {
		natives = append(natives, models.HostnameModel{
			Hostname:          types.StringValue(hostname.Hostname),
			DefaultOriginName: types.StringValue(hostname.DefaultOriginName),
			ReportCode:        types.Int64Value(int64(hostname.ReportCode)),
			TLS:               convertNativeToTLS(hostname.TLS),
			Directory:         types.StringValue(hostname.Directory),
		})
	}
	return natives
}

func convertTLSToNative(tls models.TLSModel) configuration.TLS {
	return configuration.TLS{
		NPN:                 tls.NPN.ValueBool(),
		ALPN:                tls.ALPN.ValueBool(),
		Protocols:           tls.Protocols.ValueString(),
		UseSigAlgs:          tls.UseSigAlgs.ValueBool(),
		SNI:                 tls.SNI.ValueBool(),
		SniStrict:           tls.SniStrict.ValueBool(),
		SniHostMatch:        tls.SniHostMatch.ValueBool(),
		ClientRenegotiation: tls.ClientRenegotiation.ValueBool(),
		Options:             tls.Options.ValueString(),
		CipherList:          tls.CipherList.ValueString(),
		NamedCurve:          tls.NamedCurve.ValueString(),
		OCSP:                tls.OCSP.ValueBool(),
		PEM:                 tls.PEM.ValueString(),
		CA:                  tls.CA.ValueString(),
	}
}

func convertNativeToTLS(tls configuration.TLS) models.TLSModel {
	return models.TLSModel{
		NPN:                 types.BoolValue(tls.NPN),
		ALPN:                types.BoolValue(tls.ALPN),
		Protocols:           types.StringValue(tls.Protocols),
		UseSigAlgs:          types.BoolValue(tls.UseSigAlgs),
		SNI:                 types.BoolValue(tls.SNI),
		SniStrict:           types.BoolValue(tls.SniStrict),
		SniHostMatch:        types.BoolValue(tls.SniHostMatch),
		ClientRenegotiation: types.BoolValue(tls.ClientRenegotiation),
		Options:             types.StringValue(tls.Options),
		CipherList:          types.StringValue(tls.CipherList),
		NamedCurve:          types.StringValue(tls.NamedCurve),
		OCSP:                types.BoolValue(tls.OCSP),
		PEM:                 types.StringValue(tls.PEM),
		CA:                  types.StringValue(tls.CA),
	}
}

func convertRetryToNative(retry *models.RetryModel) *configuration.Retry {
	if retry == nil {
		return nil
	}

	return &configuration.Retry{
		StatusCodes:            TypesListToIntSlice(retry.StatusCodes),
		IgnoreRetryAfterHeader: retry.IgnoreRetryAfterHeader.ValueBool(),
		AfterSeconds:           int(retry.AfterSeconds.ValueInt64()),
		MaxRequests:            int(retry.MaxRequests.ValueInt64()),
		MaxWaitSeconds:         int(retry.MaxWaitSeconds.ValueInt64()),
	}
}

func convertNativeToRetry(retry *configuration.Retry) *models.RetryModel {
	if retry == nil {
		return nil
	}

	return &models.RetryModel{
		StatusCodes:            IntSliceToTypesList(retry.StatusCodes),
		IgnoreRetryAfterHeader: types.BoolValue(retry.IgnoreRetryAfterHeader),
		AfterSeconds:           types.Int64Value(int64(retry.AfterSeconds)),
		MaxRequests:            types.Int64Value(int64(retry.MaxRequests)),
		MaxWaitSeconds:         types.Int64Value(int64(retry.MaxWaitSeconds)),
	}
}

func convertShieldsToNative(shields *models.ShieldsModel) *configuration.Shields {
	if shields == nil {
		return nil
	}

	return &configuration.Shields{
		Global: shields.Global.ValueString(),
		Apac:   shields.Apac.ValueString(),
		Emea:   shields.Emea.ValueString(),
		USWest: shields.USWest.ValueString(),
		USEast: shields.USEast.ValueString(),
	}
}

func convertNativeToShields(shields *configuration.Shields) *models.ShieldsModel {
	if shields == nil {
		return nil
	}

	return &models.ShieldsModel{
		Global: types.StringValue(shields.Global),
		Apac:   types.StringValue(shields.Apac),
		Emea:   types.StringValue(shields.Emea),
		USWest: types.StringValue(shields.USWest),
		USEast: types.StringValue(shields.USEast),
	}
}

func convertTLSVerifyToNative(tlsVerify *models.TLSVerifyModel) *configuration.TLSVerify {
	if tlsVerify == nil {
		return nil
	}

	return &configuration.TLSVerify{
		UseSNI:                   tlsVerify.UseSNI.ValueBool(),
		SNIHintAndStrictSanCheck: tlsVerify.SNIHintAndStrictSanCheck.ValueString(),
		AllowSelfSignedCerts:     tlsVerify.AllowSelfSignedCerts.ValueBool(),
		PinnedCerts:              TypesListToStringSlice(tlsVerify.PinnedCerts),
	}
}

func convertNativeToTLSVerify(tlsVerify *configuration.TLSVerify) *models.TLSVerifyModel {
	if tlsVerify == nil {
		return nil
	}

	return &models.TLSVerifyModel{
		UseSNI:                   types.BoolValue(tlsVerify.UseSNI),
		SNIHintAndStrictSanCheck: types.StringValue(tlsVerify.SNIHintAndStrictSanCheck),
		AllowSelfSignedCerts:     types.BoolValue(tlsVerify.AllowSelfSignedCerts),
		PinnedCerts:              StringSliceToTypesList(tlsVerify.PinnedCerts),
	}
}

func convertHostsToNative(hosts []models.HostModel) []configuration.Host {
	var natives []configuration.Host
	for _, host := range hosts {
		natives = append(natives, configuration.Host{
			Weight:                   int(host.Weight.ValueInt64()),
			DNSMaxTTL:                uint32(host.DNSMaxTTL.ValueInt64()),
			DNSPreference:            host.DNSPreference.ValueString(),
			MaxHardPool:              uint16(host.MaxHardPool.ValueInt64()),
			DNSMinTTL:                uint32(host.DNSMinTTL.ValueInt64()),
			Location:                 convertLocationToNative(host.Location),
			MaxPool:                  uint16(host.MaxPool.ValueInt64()),
			Balancer:                 host.Balancer.ValueString(),
			Scheme:                   host.Scheme.ValueString(),
			OverrideHostHeader:       host.OverrideHostHeader.ValueString(),
			SNIHintAndStrictSanCheck: host.SNIHintAndStrictSanCheck.ValueString(),
			UseSNI:                   host.UseSNI.ValueBool(),
		})
	}
	return natives
}

func convertNativeToHosts(hosts []configuration.Host) []models.HostModel {
	var m []models.HostModel
	for _, host := range hosts {
		m = append(m, models.HostModel{
			Weight:                   types.Int64Value(int64(host.Weight)),
			DNSMaxTTL:                types.Int64Value(int64(host.DNSMaxTTL)),
			DNSPreference:            types.StringValue(host.DNSPreference),
			MaxHardPool:              types.Int64Value(int64(host.MaxHardPool)),
			DNSMinTTL:                types.Int64Value(int64(host.DNSMinTTL)),
			Location:                 convertNativeToLocation(host.Location),
			MaxPool:                  types.Int64Value(int64(host.MaxPool)),
			Balancer:                 types.StringValue(host.Balancer),
			Scheme:                   types.StringValue(host.Scheme),
			OverrideHostHeader:       types.StringValue(host.OverrideHostHeader),
			SNIHintAndStrictSanCheck: types.StringValue(host.SNIHintAndStrictSanCheck),
			UseSNI:                   types.BoolValue(host.UseSNI),
		})
	}
	return m
}

func convertLocationToNative(locations []models.LocationModel) []configuration.Location {
	var natives []configuration.Location
	for _, location := range locations {
		natives = append(natives, configuration.Location{
			Port:     int(location.Port.ValueInt64()),
			Hostname: location.Hostname.ValueString(),
		})
	}
	return natives
}

func convertNativeToLocation(locations []configuration.Location) []models.LocationModel {
	var m []models.LocationModel
	for _, location := range locations {
		m = append(m, models.LocationModel{
			Port:     types.Int64Value(int64(location.Port)),
			Hostname: types.StringValue(location.Hostname),
		})
	}
	return m
}

func StringSliceToTypesList(slice []string) types.List {
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
		elements[i] = types.StringValue(v)
	}

	list, _ := types.ListValue(types.StringType, elements)
	return list
}

func TypesListToStringSlice(list types.List) []string {
	var result []string
	list.ElementsAs(context.Background(), &result, false)
	return result
}

func IntSliceToTypesList(slice []int) types.List {
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
		elements[i] = types.Int64Value(int64(v))
	}

	list, _ := types.ListValue(types.StringType, elements)
	return list
}

func TypesListToIntSlice(list types.List) []int {
	if list.IsNull() || list.IsUnknown() {
		return []int{}
	}

	var result []int
	list.ElementsAs(context.Background(), &result, false)
	return result
}

func StringMapToMapValue(m map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elements[k] = types.StringValue(v)
	}

	mapValue, _ := types.MapValue(types.StringType, elements)
	return mapValue
}

func MapValueToStringMap(m types.Map) map[string]string {
	elements := m.Elements()
	result := make(map[string]string, len(elements))

	for key, value := range elements {
		strValue, _ := value.(types.String)
		result[key] = strValue.ValueString()
	}

	return result
}

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNConfigurationModel struct {
	ConfigurationID        types.String    `tfsdk:"configuration_id"`
	EnvironmentID          types.String    `tfsdk:"environment_id"`
	Rules                  types.String    `tfsdk:"rules"`
	Origins                []OriginModel   `tfsdk:"origins"`
	Hostnames              []HostnameModel `tfsdk:"hostnames"`
	Experiments            types.List      `tfsdk:"experiments"`
	EdgeFunctionsSources   types.Map       `tfsdk:"edge_functions_sources"`
	EdgeFunctionInitScript types.String    `tfsdk:"edge_function_init_script"`
}

type OriginModel struct {
	Name                types.String    `tfsdk:"name"`
	Type                types.String    `tfsdk:"type"`
	Hosts               []HostModel     `tfsdk:"hosts"`
	Balancer            types.String    `tfsdk:"balancer"`
	OverrideHostHeader  types.String    `tfsdk:"override_host_header"`
	Shields             *ShieldsModel   `tfsdk:"shields"`
	PciCertifiedShields types.Bool      `tfsdk:"pci_certified_shields"`
	TLSVerify           *TLSVerifyModel `tfsdk:"tls_verify"`
	Retry               *RetryModel     `tfsdk:"retry"`
}

type HostModel struct {
	Weight                   types.Int64      `tfsdk:"weight"`
	DNSMaxTTL                types.Int64      `tfsdk:"dns_max_ttl"`
	DNSPreference            types.String     `tfsdk:"dns_preference"`
	MaxHardPool              types.Int64      `tfsdk:"max_hard_pool"`
	DNSMinTTL                types.Int64      `tfsdk:"dns_min_ttl"`
	Location                 *[]LocationModel `tfsdk:"location"`
	MaxPool                  types.Int64      `tfsdk:"max_pool"`
	Balancer                 types.String     `tfsdk:"balancer"`
	Scheme                   types.String     `tfsdk:"scheme"`
	OverrideHostHeader       types.String     `tfsdk:"override_host_header"`
	SNIHintAndStrictSanCheck types.String     `tfsdk:"sni_hint_and_strict_san_check"`
	UseSNI                   types.Bool       `tfsdk:"use_sni"`
}

type LocationModel struct {
	Port     types.Int64  `tfsdk:"port"`
	Hostname types.String `tfsdk:"hostname"`
}

type ShieldsModel struct {
	Apac   types.String `tfsdk:"apac"`
	Emea   types.String `tfsdk:"emea"`
	USWest types.String `tfsdk:"us_west"`
	USEast types.String `tfsdk:"us_east"`
}

type TLSVerifyModel struct {
	UseSNI                   types.Bool   `tfsdk:"use_sni"`
	SNIHintAndStrictSanCheck types.String `tfsdk:"sni_hint_and_strict_san_check"`
	AllowSelfSignedCerts     types.Bool   `tfsdk:"allow_self_signed_certs"`
	PinnedCerts              types.List   `tfsdk:"pinned_certs"`
}

type RetryModel struct {
	StatusCodes            types.List  `tfsdk:"status_codes"`
	IgnoreRetryAfterHeader types.Bool  `tfsdk:"ignore_retry_after_header"`
	AfterSeconds           types.Int64 `tfsdk:"after_seconds"`
	MaxRequests            types.Int64 `tfsdk:"max_requests"`
	MaxWaitSeconds         types.Int64 `tfsdk:"max_wait_seconds"`
}

type HostnameModel struct {
	Hostname          types.String `tfsdk:"hostname"`
	DefaultOriginName types.String `tfsdk:"default_origin_name"`
	ReportCode        types.Int64  `tfsdk:"report_code"`
	TLS               *TLSModel    `tfsdk:"tls"`
	Directory         types.String `tfsdk:"directory"`
}

type TLSModel struct {
	NPN                 types.Bool   `tfsdk:"npn"`
	ALPN                types.Bool   `tfsdk:"alpn"`
	Protocols           types.String `tfsdk:"protocols"`
	UseSigAlgs          types.Bool   `tfsdk:"use_sigalgs"`
	SNI                 types.Bool   `tfsdk:"sni"`
	SniStrict           types.Bool   `tfsdk:"sni_strict"`
	SniHostMatch        types.Bool   `tfsdk:"sni_host_match"`
	ClientRenegotiation types.Bool   `tfsdk:"client_renegotiation"`
	Options             types.String `tfsdk:"options"`
	CipherList          types.String `tfsdk:"cipher_list"`
	NamedCurve          types.String `tfsdk:"named_curve"`
	OCSP                types.Bool   `tfsdk:"oscp"`
	PEM                 types.String `tfsdk:"pem"`
	CA                  types.String `tfsdk:"ca"`
}

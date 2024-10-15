package dtos

import "encoding/json"

type CDNConfiguration struct {
	ConfigurationID        string            `json:"id"`
	EnvironmentID          string            `json:"environment_id"`
	Rules                  json.RawMessage   `json:"rules"`
	Origins                []Origin          `json:"origins"`
	Hostnames              []Hostname        `json:"hostnames"`
	Experiments            []string          `json:"experiments"`
	EdgeFunctionsSources   map[string]string `json:"edge_functions_sources"`
	EdgeFunctionInitScript string            `json:"edge_function_init_script"`
}

type Origin struct {
	Name                string     `json:"name"`
	Type                string     `json:"type"`
	Hosts               []Host     `json:"hosts"`
	Balancer            *string    `json:"balancer,omitempty"`
	OverrideHostHeader  *string    `json:"override_host_header,omitempty"`
	Shields             *Shields   `json:"shields,omitempty"`
	PciCertifiedShields *bool      `json:"pci_certified_shields,omitempty"`
	TLSVerify           *TLSVerify `json:"tls_verify,omitempty"`
	Retry               *Retry     `json:"retry,omitempty"`
}

type Host struct {
	Weight                   int        `json:"weight"`
	DNSMaxTTL                uint32     `json:"dns_max_ttl"`
	DNSPreference            string     `json:"dns_preference"`
	MaxHardPool              uint16     `json:"max_hard_pool"`
	DNSMinTTL                uint32     `json:"dns_min_ttl"`
	Location                 []Location `json:"location"`
	MaxPool                  uint16     `json:"max_pool"`
	Balancer                 string     `json:"balancer"`
	Scheme                   string     `json:"scheme"`
	OverrideHostHeader       string     `json:"override_host_header"`
	SNIHintAndStrictSanCheck string     `json:"sni_hint_and_strict_san_check"`
	UseSNI                   bool       `json:"use_sni"`
}

type Location struct {
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
}

type Shields struct {
	Apac   string `json:"apac"`
	Emea   string `json:"emea"`
	USWest string `json:"us_west"`
	USEast string `json:"us_east"`
}

type TLSVerify struct {
	UseSNI                   bool     `json:"use_sni"`
	SNIHintAndStrictSanCheck string   `json:"sni_hint_and_strict_san_check"`
	AllowSelfSignedCerts     bool     `json:"allow_self_signed_certs"`
	PinnedCerts              []string `json:"pinned_certs"`
}

type Retry struct {
	StatusCodes            []int `json:"status_codes"`
	IgnoreRetryAfterHeader bool  `json:"ignore_retry_after_header"`
	AfterSeconds           int   `json:"after_seconds"`
	MaxRequests            int   `json:"max_requests"`
	MaxWaitSeconds         int   `json:"max_wait_seconds"`
}

type Hostname struct {
	Hostname          string `json:"hostname"`
	DefaultOriginName string `json:"default_origin_name"`
	ReportCode        int    `json:"report_code"`
	TLS               TLS    `json:"tls"`
	Directory         string `json:"directory"`
}

type TLS struct {
	NPN                 bool   `json:"npn"`
	ALPN                bool   `json:"alpn"`
	Protocols           string `json:"protocols"`
	UseSigAlgs          bool   `json:"use_sigalgs"`
	SNI                 bool   `json:"sni"`
	SniStrict           bool   `json:"sni_strict"`
	SniHostMatch        bool   `json:"sni_host_match"`
	ClientRenegotiation bool   `json:"client_renegotiation"`
	Options             string `json:"options"`
	CipherList          string `json:"cipher_list"`
	NamedCurve          string `json:"named_curve"`
	OCSP                bool   `json:"oscp"`
	PEM                 string `json:"pem"`
	CA                  string `json:"ca"`
}

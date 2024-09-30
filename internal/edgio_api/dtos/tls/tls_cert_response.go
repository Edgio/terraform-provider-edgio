package tls

type TLSCertResponse struct {
	ID               string   `json:"id"`
	EnvironmentID    string   `json:"environment_id"`
	PrimaryCert      string   `json:"primary_cert"`
	IntermediateCert string   `json:"intermediate_cert"`
	Expiration       string   `json:"expiration"`
	Status           string   `json:"status"`
	Generated        bool     `json:"generated"`
	Serial           string   `json:"serial"`
	CommonName       string   `json:"common_name"`
	AlternativeNames []string `json:"alternative_names"`
	ActivationError  string   `json:"activation_error"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

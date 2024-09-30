package tls

type UploadTlsCertRequest struct {
	EnvironmentID    string `json:"environment_id"`
	PrimaryCert      string `json:"primary_cert"`
	IntermediateCert string `json:"intermediate_cert"`
	PrivateKey       string `json:"private_key"`
}

package tls

type TLSCertSResponse struct {
	EnvironmentID string            `json:"environment_id"`
	TotalItems    int32             `json:"total_items"`
	Certificates  []TLSCertResponse `json:"items"`
}

package purge

type PurgeRequest struct {
	EnvironmentID string   `json:"environment_id"`
	PurgeType     string   `json:"purge_type"`
	Values        []string `json:"values"`
	Hostname      *string  `json:"hostname"`
}

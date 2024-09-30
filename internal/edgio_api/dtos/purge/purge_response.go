package purge

import "time"

// API Response structure
type PurgeResponse struct {
	ID                 string    `json:"id"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	CompletedAt        time.Time `json:"completed_at"`
	ProgressPercentage float32   `json:"progress_percentage"`
}

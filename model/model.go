package model

import "time"

// URL entity structure
type URL struct {
	ID        string     `json:"id" bson:"_id"`
	LongURL   string     `json:"longurl"`
	CreatedAt *time.Time `json:"createdAt"`
	CreatedBy string     `json:"createdBy"`
	UpdatedAt *time.Time `json:"updatedAt"`
	UpdatedBy string     `json:"updatedBy"`
	Metrics   struct {
		RoutedCount int64 `json:"routedCount"`
	} `json:"metrics"`
}

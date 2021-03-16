package model

import "time"

// Tag -
type Tag struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

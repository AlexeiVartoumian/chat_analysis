package models

import "time"

type JobLifeCycle struct {
	JobId             int        `json:"job_id,omitempty" db:"job_id"`
	JobState          string     `json:"job_state,omitempty" db:"job_state"`
	FirstSeenAt       time.Time  `json:"first_seen_at,omitempty" db:"first_seen_at"`
	LastSeenListedAt  time.Time  `json:"last_seen_listed_at,omitempty" db:"last_seen_listed_at"`
	FirstSeenClosedAt *time.Time `json:"first_seen_closed_at,omitempty" db:"first_seen_closed_at"`
	NextScanAt        *time.Time `json:"next_scan_at,omitempty" db:"next_scan_at"`
	SuspendedCount    int        `json:"suspended_count,omitempty" db:"suspended_count"`
}

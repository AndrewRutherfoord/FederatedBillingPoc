package models

import "time"

// TimeRangeQuery is the standard `from`/`to` querystring pair used by
// time-bounded list endpoints across the billing system's HTTP APIs.
type TimeRangeQuery struct {
	From time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To   time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
}

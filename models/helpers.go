package models

import "time"

type PeriodicFunctions struct {
	Fn       func()
	Interval time.Duration
}

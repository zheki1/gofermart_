package models

import "time"

type Withdrawal struct {
	Order       string
	Sum         float64
	ProcessedAt time.Time
}

type Balance struct {
	Current   float64
	Withdrawn float64
}

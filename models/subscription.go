package models

import "time"

type Subscription struct {
	ID              int
	UserID          int
	CategoryID      int
	SearchText      string
	PriceMin        int
	PriceMax        int
	RegionID        int
	ProTypes        string // "0,3,5"
	FilterSignature string
	CreatedAt       time.Time
}

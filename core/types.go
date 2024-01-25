package core

type PriceType string
type RecuringInterval string

const (
	PriceTypeSingle       PriceType        = "single"
	PriceTypeRecurring    PriceType        = "recurring"
	RecuringIntervalDay   RecuringInterval = "day"
	RecuringIntervalWeek  RecuringInterval = "week"
	RecuringIntervalMonth RecuringInterval = "month"
	RecuringIntervalYear  RecuringInterval = "year"
)
package models

type Subscription struct {
	ID                uint64  `gorm:"primaryKey"`
	Amount            float64 `gorm:"not null"`
	Active            bool    `gorm:"index;not null;default:TRUE"`
	Interval          int     `gorm:"not null"`
	Currency          string  `gorm:"not null"`
	AccountAddress    string  `gorm:"index;not null"`
	DestinationAdress string  `gorm:"index;not null"`
}

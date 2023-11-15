package models

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID                uuid.UUID `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EntityId          uuid.UUID // This is a unique id of the entity that the webhook is intented for eg. A Merchant Id
	HookKey           string
	Payload           string
	Status            WebhookStatus    `gorm:"index;default:Created"`
	Processor         WebhookProcessor `gorm:"default:Convoy"`
	ProcessorResponse string
	FailureReason     string
}

type WebhookStatus string
type WebhookProcessor string

const (
	WebhookSent             WebhookStatus    = "Sent"
	WebhookCreated          WebhookStatus    = "Created"
	WebhookFailedToSend     WebhookStatus    = "FailedToSend"
	ConvoyWebhookProccessor WebhookProcessor = "Convoy"
)

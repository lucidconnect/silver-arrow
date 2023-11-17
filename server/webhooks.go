package server

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
)

const (
	PaymentSuccessWebHook WebHookType = "Payment.Success"
	PaymentFailedWebHook  WebHookType = "Payment.Failed"
	PaymentPendingWebHook WebHookType = "Payment.Pending"
)

type WebHookType string

type PaymentWebHookBody struct {
	Event string         `json:"event"`
	Data  WebHookPayment `json:"data"`
}

type WebHookPayment struct {
	Status           string  `json:"status"`
	Amount           float64 `json:"amount"`
	Source           string  `json:"source"`
	Destination      string  `json:"destination"`
	PaymentReference string  `json:"paymentReference"`
}

func (s *Server) TriggerWebhook(merchant models.Merchant, paymentReference string) {
	var event string

	// ws := wallet.NewWalletService(s.database, nil, 0)
	// payment, err := ws.FetchPayment(paymentReference)
	// if err != nil {
	// 	log.Err(err).Send()
	// }
	ref, err := uuid.Parse(paymentReference)
	if err != nil {
		log.Err(err).Msg("invalid reference")
		// return nil, errors.New("invalid reference")
		return
	}
	payment, err := s.database.FindPaymentByReference(ref)
	if err != nil {
		log.Err(err).Msgf("payment [%v] does not exist", paymentReference)
		// return nil, errors.New("invalid payment reference")
		return
	}

	switch payment.Status {
	case models.PaymentStatusSuccess:
		event = string(PaymentSuccessWebHook)
	case models.PaymentStatusFailed:
		event = string(PaymentFailedWebHook)
	case models.PaymentStatusPending:
		event = string(PaymentPendingWebHook)
	}

	amount := parseTransferAmountFloat(payment.Token, payment.Amount)
	requestBody := PaymentWebHookBody{
		Event: event,
		Data: WebHookPayment{
			Status: string(payment.Status),
			Amount: amount,
			Source: payment.Source,
			// Destination: payment.Destination,
			PaymentReference: payment.Reference.String(),
		},
	}

	stringifyBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	id := uuid.New()
	webhookEvent := models.WebhookEvent{
		ID:        id,
		EntityId:  merchant.ID,
		Payload:   string(stringifyBody),
		Status:    models.WebhookCreated,
		Processor: models.ConvoyWebhookProccessor,
		// HookKey:   merchant.WebHookToken,
	}

	// create webhook event
	s.database.CreateWebhookEvent(&webhookEvent)
	// create webhook event on convoy
	convoyClient := convoy.New(convoy.Options{
		APIKey:    os.Getenv("CONVOY_API_KEY"),
		ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	eventResponse, err := convoyClient.Events.Create(&convoy.CreateEventRequest{
		EventType:  event,
		Data:       stringifyBody,
		EndpointID: merchant.ConvoyEndpointID,
	}, &convoy.EventQueryParam{
		GroupID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	if err != nil {
		webhookEvent.FailureReason = err.Error()
		webhookEvent.Status = models.WebhookFailedToSend
		// log.Println(payment.ID, "Failed to post webhook event to Convoy because of ", err.Error())
	}

	// we keep a Serialized format of the callback responses for testing/mocking purposes
	if serializedResponse, err := json.Marshal(eventResponse); err != nil {
		webhookEvent.ProcessorResponse = "Failed to marshal webhook response : " + err.Error()
	} else {
		webhookEvent.ProcessorResponse = string(serializedResponse)
	}

	webhookEvent.Status = models.WebhookSent
	// if err := utils.DB.Save(&webhookEvent).Error; err != nil {
	// 	formatedLog := fmt.Sprintf("Payment (%s) failed to persit updated webhook event because of %s", payment.ID, err.Error())
	// 	webhookEvent.FailureReason += formatedLog
	// 	log.Println(formatedLog)
	// }

	if err := s.database.UpdateWebhookEvent(&webhookEvent); err != nil {
		formatedLog := fmt.Sprintf("Payment (%s) failed to persit updated webhook event because of %s", payment.ID, err.Error())
		webhookEvent.FailureReason += formatedLog
		log.Err(err).Msg(formatedLog)
	}

	paymentUpdate := map[string]any{
		"acknowledged":            true,
		"webhook_acknowledged_at": time.Now(),
	}
	if err := s.database.UpdatePayment(payment.ID, paymentUpdate); err != nil {
		formatedMessage := fmt.Sprintf("Payment (%s) failed to update but a webhook was forwared to the merchant because of %s", payment.ID, err.Error())
		webhookEvent.FailureReason += formatedMessage
		log.Err(err).Msg(formatedMessage)
	}
}

func parseTransferAmountFloat(token string, amount int64) float64 {
	var divisor int
	if token == "USDC" || token == "USDT" {
		divisor = 6
	} else {
		divisor = 18
	}
	minorFactor := math.Pow10(divisor)
	parsedAmount := float64(amount) / minorFactor

	return parsedAmount
}

/**

func TriggerWebHook(payment models.Payments, merchant models.Merchants, userid string) {


	webhookEvent := models.WebhookEvent{
		EntityId:  merchant.ID,
		Payload:   string(stringifyBody),
		Status:    models.WebhookCreated,
		Processor: models.ConvoyWebhookProccessor,
		HookKey:   merchant.WebHookToken,
	}

	if err := utils.DB.Create(&webhookEvent).Error; err != nil {
		log.Println(payment.ID, "Error creating webhook event: ", err)
		return
	}

	convoyClient := convoy.New(convoy.Options{
		APIKey: os.Getenv("CONVOY_API_KEY"),
		ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	eventResponse, err := convoyClient.Events.Create(&convoy.CreateEventRequest{
		EventType:  event,
		Data:       stringifyBody,
		EndpointID: merchant.EndpointID,
	}, &convoy.EventQueryParam{
		GroupID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	if err != nil {
		webhookEvent.FailureReason = err.Error()
		webhookEvent.Status = models.WebhookFailedToSend
		log.Println(payment.ID, "Failed to post webhook event to Convoy because of ", err.Error())
	}

	// we keep a Serialized format of the callback responses for testing/mocking purposes
	if serializedResponse, err := json.Marshal(eventResponse); err != nil {
		webhookEvent.ProcessorResponse = "Failed to marshal webhook response : " + err.Error()
	} else {
		webhookEvent.ProcessorResponse = string(serializedResponse)
	}

	webhookEvent.Status = models.WebhookSent
	if err := utils.DB.Save(&webhookEvent).Error; err != nil {
		formatedLog := fmt.Sprintf("Payment (%s) failed to persit updated webhook event because of %s", payment.ID, err.Error())
		webhookEvent.FailureReason += formatedLog
		log.Println(formatedLog)
	}

	payment.Acknowledged = true
	payment.WebhookAcknowledgedAt = time.Now()
	if err := utils.DB.Save(&payment).Error; err != nil {
		formatedMessage := fmt.Sprintf("Payment (%s) failed to update but a webhook was forwared to the merchant because of %s", payment.ID, err.Error())
		webhookEvent.FailureReason += formatedMessage
		log.Println(formatedMessage)
	}
}

*/

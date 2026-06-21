package cloudserviceprovider

import (
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

type GetBillingAccountRecordsRequest struct {
	BillingAccountID string    `json:"billing_account_id"`
	From             time.Time `json:"from"`
	To               time.Time `json:"to"`
}

type GetBillingAccountRecordsResponse struct {
	Batches []models.ChargeBatchDetail `json:"batches"`
	Count   int                        `json:"count"`
	From    time.Time                  `json:"from"`
	To      time.Time                  `json:"to"`
}

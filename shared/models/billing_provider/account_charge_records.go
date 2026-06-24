package billingprovider

import (
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

type GetBillingAccountRecordsResponse struct {
	Batches []models.ChargeBatch `json:"batches"`
	Count   int                  `json:"count"`
	From    time.Time            `json:"from"`
	To      time.Time            `json:"to"`
}

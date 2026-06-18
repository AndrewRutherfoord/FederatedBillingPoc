package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	billingprovidermodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	ReturnURL string `json:"return_url" binding:"required"`
}

type registerResponse struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirect_url"`
}

func (s *Server) RegisterBillingAccount(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := s.repos.BillingAccounts.Create(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create billing account"})
		return
	}

	baseURL := fmt.Sprintf("%s://%s", scheme(c.Request), c.Request.Host)
	redirectURL := fmt.Sprintf("%s/billing/accounts/%s/onboard?return_url=%s", baseURL, account.ID, req.ReturnURL)

	c.JSON(http.StatusCreated, registerResponse{
		ID:          account.ID,
		RedirectURL: redirectURL,
	})
}

func (s *Server) OnboardForm(c *gin.Context) {
	id := c.Param("id")
	returnURL := c.Query("return_url")

	c.HTML(http.StatusOK, "onboard.html", gin.H{
		"AccountID": id,
		"ReturnURL": returnURL,
		"Name":      "",
		"Email":     "",
		"Error":     "",
	})
}

func (s *Server) OnboardSubmit(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	returnURL := c.PostForm("return_url")

	if name == "" || email == "" {
		c.HTML(http.StatusUnprocessableEntity, "onboard.html", gin.H{
			"AccountID": id,
			"ReturnURL": returnURL,
			"Name":      name,
			"Email":     email,
			"Error":     "Name and email are required.",
		})
		return
	}

	if _, err := s.repos.BillingAccounts.Update(c.Request.Context(), id, name, email); err != nil {
		c.HTML(http.StatusInternalServerError, "onboard.html", gin.H{
			"AccountID": id,
			"ReturnURL": returnURL,
			"Name":      name,
			"Email":     email,
			"Error":     "Something went wrong. Please try again.",
		})
		return
	}

	u, err := url.Parse(returnURL)
	if err != nil {
		c.Redirect(http.StatusSeeOther, returnURL)
		return
	}
	q := u.Query()
	q.Set("account_id", id)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusSeeOther, u.String())
}

func (s *Server) GetBillingAccountRecords(c *gin.Context) {
	var req billingprovidermodels.GetBillingAccountRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batches, err := s.repos.CostBatch.GetByBillingAccountAndTimeRange(c.Request.Context(), req.BillingAccountID, req.From, req.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch billing records"})
		return
	}

	records := make([]models.AggregatedChargeRecord, 0, len(batches))
	for _, batch := range batches {
		records = append(records, models.AggregatedChargeRecord{
			BillingRecord: models.BillingRecord{
				BillingProviderID:  s.repos.Provider.ID,
				ResourceProviderID: batch.CloudServiceProviderID,
				BillingAccountID:   batch.BillingAccountID,
			},
			BatchID:         batch.ID,
			TotalBilledCost: batch.TotalCost,
			BilledCurrency:  "EUR",
			LineItemCount:   batch.TotalItems,
			BatchHash:       batch.MerkelRoot,
			BatchSignature:  batch.Signature,
			CreatedAt:       batch.CreatedAt,
		})
	}

	response := billingprovidermodels.GetBillingAccountRecordsResponse{
		Records: records,
		Count:   len(records),
		From:    req.From,
		To:      req.To,
	}

	c.JSON(http.StatusOK, response)
}

func scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "http"
}

package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
	"github.com/gin-gonic/gin"
)

// InitiateBillingAccountOnboarding godoc
//
//	@Summary	Initiate a CSP billing account onboarding session
//	@Description	Called by the customer-billing service. Creates a pending onboarding session and returns a redirect URL for the customer to complete setup.
//	@Tags		billing-accounts
//	@Accept		json
//	@Produce	json
//	@Param		body	body		cspsharedmodels.RegisterLinkedCloudProviderRequest	true	"Registration request"
//	@Success	201		{object}	cspsharedmodels.RegisterLinkedCloudProviderResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/billing/accounts [post]
func (s *Server) InitiateBillingAccountOnboarding(c *gin.Context) {
	var req cspsharedmodels.RegisterLinkedCloudProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := s.repos.OnboardingSessions.Create(c.Request.Context(), req.BillingProviderID, req.BillingAccountID, req.ReturnURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create onboarding session"})
		return
	}

	baseURL := fmt.Sprintf("%s://%s", requestScheme(c.Request), c.Request.Host)
	redirectURL := fmt.Sprintf("%s/onboarding/%s", baseURL, session.ID)

	c.JSON(http.StatusCreated, cspsharedmodels.RegisterLinkedCloudProviderResponse{
		ID:          session.ID,
		RedirectURL: redirectURL,
	})
}

// OnboardingForm godoc
//
//	@Summary	Render the CSP account onboarding form
//	@Tags		onboarding
//	@Produce	html
//	@Param		session_id	path	string	true	"Onboarding session ID"
//	@Success	200
//	@Failure	404	{object}	map[string]string
//	@Router		/onboarding/{session_id} [get]
func (s *Server) OnboardingForm(c *gin.Context) {
	sessionID := c.Param("session_id")
	session, err := s.repos.OnboardingSessions.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "onboarding session not found"})
		return
	}

	c.HTML(http.StatusOK, "csp_onboard.html", gin.H{
		"SessionID":         session.ID,
		"BillingAccountID":  session.BillingAccountID,
		"BillingProviderID": session.BillingProviderID,
		"ProviderName":      s.config.ProviderName,
		"Mode":              "new",
		"CustomerID":        "",
		"Name":              "",
		"Email":             "",
		"Error":             "",
	})
}

// OnboardingSubmit godoc
//
//	@Summary	Submit the CSP account onboarding form
//	@Tags		onboarding
//	@Accept		application/x-www-form-urlencoded
//	@Param		session_id	path	string	true	"Onboarding session ID"
//	@Success	303
//	@Failure	404	{object}	map[string]string
//	@Router		/onboarding/{session_id} [post]
func (s *Server) OnboardingSubmit(c *gin.Context) {
	sessionID := c.Param("session_id")
	session, err := s.repos.OnboardingSessions.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "onboarding session not found"})
		return
	}

	mode := c.PostForm("mode")

	renderError := func(msg string, customerID, name, email string) {
		c.HTML(http.StatusUnprocessableEntity, "csp_onboard.html", gin.H{
			"SessionID":         session.ID,
			"BillingAccountID":  session.BillingAccountID,
			"BillingProviderID": session.BillingProviderID,
			"ProviderName":      s.config.ProviderName,
			"Mode":              mode,
			"CustomerID":        customerID,
			"Name":              name,
			"Email":             email,
			"Error":             msg,
		})
	}

	var customerID string

	if mode == "existing" {
		customerID = c.PostForm("customer_id")
		if customerID == "" {
			renderError("Customer ID is required.", customerID, "", "")
			return
		}
		if _, err := s.repos.Customers.GetByID(c.Request.Context(), customerID); err != nil {
			renderError("No account found with that Customer ID.", customerID, "", "")
			return
		}
	} else {
		name := c.PostForm("name")
		email := c.PostForm("email")
		if name == "" || email == "" {
			renderError("Name and email are required.", "", name, email)
			return
		}
		customer, err := s.repos.Customers.Create(c.Request.Context(), name, email)
		if err != nil {
			renderError("Failed to create account. That email may already be registered.", "", name, email)
			return
		}
		customerID = customer.ID
	}

	// Link the billing provider account to this CSP customer.
	_, err = s.repos.BillingAccounts.Create(c.Request.Context(), db.BillingAccount{
		AccountID:         session.BillingAccountID,
		BillingProviderID: session.BillingProviderID,
		CustomerID:        customerID,
	})
	if err != nil {
		renderError("Failed to link billing account. It may already be linked.", "", "", "")
		return
	}

	u, _ := url.Parse(session.ReturnURL)
	q := u.Query()
	q.Set("account_id", session.BillingAccountID)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusSeeOther, u.String())
}

func requestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "http"
}

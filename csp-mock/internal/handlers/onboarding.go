package handlers

import (
	"context"
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

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/onboarding/%s/complete", session.ID))
}

// OnboardingComplete godoc
//
//	@Summary	Show the linked account info after onboarding
//	@Tags		onboarding
//	@Produce	html
//	@Param		session_id	path	string	true	"Onboarding session ID"
//	@Success	200
//	@Failure	404	{object}	map[string]string
//	@Router		/onboarding/{session_id}/complete [get]
func (s *Server) OnboardingComplete(c *gin.Context) {
	sessionID := c.Param("session_id")
	session, err := s.repos.OnboardingSessions.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "onboarding session not found"})
		return
	}

	info, err := s.onboardingAccountInfo(c.Request.Context(), session)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "linked billing account not found"})
		return
	}

	returnURL, err := returnURLWithAccountID(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid return url"})
		return
	}

	c.HTML(http.StatusOK, "csp_onboard_complete.html", gin.H{
		"ProviderName": s.config.ProviderName,
		"Info":         info,
		"ReturnURL":    returnURL,
		"DownloadURL":  fmt.Sprintf("/onboarding/%s/complete/download", session.ID),
	})
}

// OnboardingCompleteDownload godoc
//
//	@Summary	Download the linked account info as JSON
//	@Tags		onboarding
//	@Produce	json
//	@Param		session_id	path	string	true	"Onboarding session ID"
//	@Success	200	{object}	OnboardingAccountInfo
//	@Failure	404	{object}	map[string]string
//	@Router		/onboarding/{session_id}/complete/download [get]
func (s *Server) OnboardingCompleteDownload(c *gin.Context) {
	sessionID := c.Param("session_id")
	session, err := s.repos.OnboardingSessions.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "onboarding session not found"})
		return
	}

	info, err := s.onboardingAccountInfo(c.Request.Context(), session)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "linked billing account not found"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s_%s.json"`, info.ProviderID, info.BillingAccountID))
	c.JSON(http.StatusOK, info)
}

// OnboardingAccountInfo is the account information shown/downloaded once onboarding completes.
type OnboardingAccountInfo struct {
	CustomerID        string `json:"customer_id"`
	CustomerName      string `json:"customer_name"`
	CustomerEmail     string `json:"customer_email"`
	BillingAccountID  string `json:"billing_account_id"`
	BillingProviderID string `json:"billing_provider_id"`
	ProviderID        string `json:"provider_id"`
	ProviderName      string `json:"provider_name"`
	Host              string `json:"host"`
}

func (s *Server) onboardingAccountInfo(ctx context.Context, session *db.OnboardingSession) (*OnboardingAccountInfo, error) {
	account, err := s.repos.BillingAccounts.GetByAccountID(ctx, session.BillingAccountID)
	if err != nil {
		return nil, err
	}
	customer, err := s.repos.Customers.GetByID(ctx, account.CustomerID)
	if err != nil {
		return nil, err
	}
	return &OnboardingAccountInfo{
		CustomerID:        customer.ID,
		CustomerName:      customer.Name,
		CustomerEmail:     customer.Email,
		BillingAccountID:  account.AccountID,
		BillingProviderID: account.BillingProviderID,
		ProviderID:        s.config.ProviderID,
		ProviderName:      s.config.ProviderName,
		Host:              s.config.CustomerAPIEndpointURL,
	}, nil
}

func returnURLWithAccountID(session *db.OnboardingSession) (string, error) {
	u, err := url.Parse(session.ReturnURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("account_id", session.BillingAccountID)
	u.RawQuery = q.Encode()
	return u.String(), nil
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

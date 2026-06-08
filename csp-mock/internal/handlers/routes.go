package handlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", s.Health)
	r.GET("/clock/current", s.GetCurrentTime)
	r.POST("/clock/advance", s.AdvanceTime)

	r.GET("/resource-types", s.ListResourceTypes)
	r.GET("/resource-types/:id", s.GetResourceType)

	r.POST("/customer/register", s.RegisterCustomer)

	// Routes below require a valid customer in the Authorization header.
	authed := r.Group("/", middleware.Auth(s.repos.Customers))
	{
		authed.GET("/customer", s.GetCustomer)

		authed.GET("/resources", s.ListResources)
		authed.POST("/resources", s.CreateResource)
		authed.DELETE("/resources/:id", s.DeleteResource)

		authed.GET("/billing-accounts", s.ListBillingAccounts)
		authed.GET("/billing-accounts/:account_id", s.GetBillingAccount)
		authed.POST("/billing-accounts", s.CreateBillingAccount)
	}
}

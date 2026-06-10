package customerhandlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
	"github.com/gin-gonic/gin"
)

type CustomerServer struct {
	repos     *repository.Repos
	clock     util.Clock
	scheduler *scheduler.Scheduler
}

func NewCustomerServer(repos *repository.Repos, clock util.Clock) *CustomerServer {
	return &CustomerServer{repos: repos, clock: clock}
}

func (cs *CustomerServer) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/customer")

	group.GET("/resource-types", cs.ListResourceTypes)
	group.GET("/resource-types/:id", cs.GetResourceType)

	group.POST("/customer/register", cs.RegisterCustomer)

	// Routes below require a valid customer in the Authorization header.
	authed := group.Group("/", middleware.Auth(cs.repos.Customers))
	{
		authed.GET("/customer", cs.GetCustomer)

		authed.GET("/resources", cs.ListResources)
		authed.POST("/resources", cs.CreateResource)
		authed.DELETE("/resources/:id", cs.DeleteResource)

		authed.GET("/billing-accounts", cs.ListBillingAccounts)
		authed.GET("/billing-accounts/:account_id", cs.GetBillingAccount)
		authed.POST("/billing-accounts", cs.CreateBillingAccount)
	}
}

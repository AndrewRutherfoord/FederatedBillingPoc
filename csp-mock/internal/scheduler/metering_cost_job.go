package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared-models"
	"github.com/shopspring/decimal"
)

type RecordMeteringAndCostJob struct {
	id     string
	repos  *repository.Repos
	config *config.CspConfig
}

func NewRecordMeteringAndCostJob(id string, repos *repository.Repos, config *config.CspConfig) *RecordMeteringAndCostJob {
	return &RecordMeteringAndCostJob{id: id, repos: repos, config: config}
}

func (j *RecordMeteringAndCostJob) ID() string {
	return j.id
}

func (j *RecordMeteringAndCostJob) createBaseCostRecord(ctx context.Context, resourceType *config.ResourceType) string {
	// Implementation for creating base cost record
	return ""
}

func (j *RecordMeteringAndCostJob) createPerHourCostRecord(baseBuilder *sharedmodels.FocusLineItemBuilder, resourceType *config.ResourceType, resource *db.Resource) *sharedmodels.FocusLineItem {
	// TODO: Currently assumes this is run on the hour every hour. Needs to actually calculate the hours since last run.
	qty := decimal.NewFromInt(1)
	hour := "hour"
	resourceID := resource.ID
	resourceTypeStr := resource.ResourceType

	return baseBuilder.Copy().
		WithServiceAndCosts(
			sharedmodels.ServiceCategoryCompute,
			sharedmodels.ServiceSubcategoryVirtualMachines,
			resourceType.ID,
			resourceType.Pricing.UnitPrice,
			resourceType.Pricing.UnitPrice,
			resourceType.Pricing.UnitPrice,
			resourceType.Pricing.UnitPrice,
		).
		WithConsumption(&qty, &hour).
		WithResource(&resourceID, nil, &resourceTypeStr).
		Build()
}

func (j *RecordMeteringAndCostJob) createPerGbHourCostRecord(baseBuilder *sharedmodels.FocusLineItemBuilder, resourceType *config.ResourceType, resource *db.Resource) *sharedmodels.FocusLineItem {
	// TODO: Currently assumes this is run on the hour every hour. Needs to actually calculate the hours since last run.
	if resource.StorageGB == nil {
		log.Printf("Resource %s has no storage GB value for per-GB-hour billing", resource.ID)
		return nil
	}

	gbHour := "GB-hour"
	billedCost := resource.StorageGB.Mul(resourceType.Pricing.UnitPrice)
	resourceID := resource.ID
	resourceTypeStr := resource.ResourceType

	return baseBuilder.Copy().
		WithServiceAndCosts(
			sharedmodels.ServiceCategoryStorage,
			sharedmodels.ServiceSubcategoryBlockStorage,
			resourceType.ID,
			billedCost,
			billedCost,
			billedCost,
			billedCost,
		).
		WithConsumption(resource.StorageGB, &gbHour).
		WithResource(&resourceID, nil, &resourceTypeStr).
		Build()
}

func (j *RecordMeteringAndCostJob) createResourceCostRecord(ctx context.Context, baseBuilder *sharedmodels.FocusLineItemBuilder, billingAccount *db.BillingAccount, resource *db.Resource) *sharedmodels.FocusLineItem {
	resourceType, err := j.repos.ResourceTypes.GetById(ctx, resource.ResourceType)
	if err != nil {
		log.Printf("Error fetching resource type '%s' for resource '%s': %v", resource.ResourceType, resource.ID, err)
		return nil
	}

	switch resourceType.Pricing.Model {
	case "per_hour":
		return j.createPerHourCostRecord(baseBuilder, resourceType, resource)
	case "per_gb_hour":
		return j.createPerGbHourCostRecord(baseBuilder, resourceType, resource)
	default:
		log.Printf("Unknown pricing model '%s' for resource type '%s'", resourceType.Pricing.Model, resourceType.ID)
		return nil
	}
}

func (j *RecordMeteringAndCostJob) processCustomerResources(ctx context.Context, startTime time.Time, billingAccount db.BillingAccount) error {
	resources, err := j.repos.Resources.ListByBillingAccountID(ctx, billingAccount.AccountID)
	if err != nil {
		log.Printf("Error fetching resources for billing account %s: %v", billingAccount.AccountID, err)
		return err
	}

	// TODO: Do this properly based on BP billing period, not just current month
	billingPeriodStart := time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingPeriodEnd := billingPeriodStart.AddDate(0, 1, 0)
	chargePeriodStart := startTime.Add(-1 * time.Hour)
	chargePeriodEnd := startTime

	baseBuilder := sharedmodels.NewFocusLineItemBuilder().
		WithProviderAndPeriod(
			billingAccount.AccountID,
			"Customer", // TODO: Not really sure what to put here.
			j.config.Currency,
			billingPeriodStart,
			billingPeriodEnd,
			"usage",
			"Usage-Based",
			chargePeriodStart,
			chargePeriodEnd,
			j.config.ProviderName,
			j.config.ProviderName,
			j.config.ProviderName, // TODO: This should probably be the BP?
		).
		WithLocation(
			j.config.AvailabilityZone,
			j.config.RegionName,
			j.config.RegionID,
		)

	for _, resource := range resources {
		lineItem := j.createResourceCostRecord(ctx, baseBuilder, &billingAccount, resource)
		if lineItem == nil {
			log.Printf("Failed to create cost record for resource %s", resource.ID)
			continue
		}

		// Save line item to repository
		if err := j.repos.Focus.Insert(ctx, *lineItem); err != nil {
			log.Printf("Error saving line item for resource %s: %v", resource.ID, err)
			continue
		}
	}
	return nil
}

func (j *RecordMeteringAndCostJob) Execute(ctx context.Context, startTime time.Time) error {
	log.Printf("RecordMeteringAndCostJob executed: %s", j.id)

	billingAccount, err := j.repos.BillingAccounts.List(ctx)
	if err != nil {
		log.Printf("Error occurred while fetching customers: %v", err)
		return err
	}

	for _, ba := range billingAccount {
		j.processCustomerResources(ctx, startTime, ba)
	}

	return nil
}

package scheduler

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/port"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/shopspring/decimal"
)

type RecordMeteringAndCostJob struct {
	id           string
	repos        *repository.Repos
	config       *config.CspConfig
	batchMaxSize int
	sender       port.BPSender
}

func NewRecordMeteringAndCostJob(id string, repos *repository.Repos, config *config.CspConfig, sender port.BPSender) *RecordMeteringAndCostJob {
	return &RecordMeteringAndCostJob{id: id, repos: repos, config: config, batchMaxSize: 1000, sender: sender}
}

func (j *RecordMeteringAndCostJob) ID() string {
	return j.id
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

func (j *RecordMeteringAndCostJob) createResourceCostRecord(ctx context.Context, baseBuilder *sharedmodels.FocusLineItemBuilder, resource *db.Resource) *sharedmodels.FocusLineItem {
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

func (j *RecordMeteringAndCostJob) processBillingAccountResourcesBatch(ctx context.Context, startTime time.Time, billingAccount db.BillingAccount, resources []*db.Resource) error {

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

	lineItems := make([]sharedmodels.FocusLineItem, 0, len(resources))
	var hashes []string

	for _, resource := range resources {
		lineItem := j.createResourceCostRecord(ctx, baseBuilder, resource)
		if lineItem == nil {
			log.Printf("Failed to create cost record for resource %s", resource.ID)
			continue
		}

		hash, err := shared.CanonicalHash(lineItem)
		if err != nil {
			log.Printf("Error hashing line item for resource %s: %v", resource.ID, err)
			continue
		}

		hashes = append(hashes, hash)
		lineItems = append(lineItems, *lineItem)
	}

	// Build merkle tree and create batch
	// Sort hashes for deterministic merkle tree root regardless of query order
	sort.Strings(hashes)
	merkleTree := shared.NewMerkleTree(hashes)
	var totalCost float64
	for _, item := range lineItems {
		totalCost += item.BilledCost.InexactFloat64()
	}

	batch, err := j.repos.CostBatch.Create(ctx, billingAccount.AccountID, billingAccount.BillingProviderID, j.config.ProviderID, merkleTree.Root(), len(lineItems), totalCost, startTime)
	if err != nil {
		log.Printf("Error creating cost batch for billing account %s: %v", billingAccount.AccountID, err)
		return err
	}

	// Save line items with batch ID
	if err := j.repos.Focus.InsertBatchWithBatchID(ctx, lineItems, batch.ID); err != nil {
		log.Printf("Error saving line items for batch %s: %v", batch.ID, err)
		return err
	}

	// Send it to the billing provider
	return j.sender.SendAggregatedChargeRecord(ctx, sharedmodels.AggregatedChargeRecord{
		BillingRecord: sharedmodels.BillingRecord{
			BillingProviderID:  batch.BillingProviderID,
			ResourceProviderID: batch.ResourceProviderID,
			BillingAccountID:   batch.BillingAccountID,
		},
		BatchID:         batch.ID,
		TotalBilledCost: batch.TotalCost,
		BilledCurrency:  batch.BilledCurrency,
		LineItemCount:   batch.TotalItems,
		BatchHash:       batch.MerkelRoot,
		BatchSignature:  "", // TODO: Sign the batch with the CSP's private key
		CreatedAt:       batch.CreatedAt,
	})
}

func (j *RecordMeteringAndCostJob) Execute(ctx context.Context, startTime time.Time) error {
	log.Printf("RecordMeteringAndCostJob executed: %s", j.id)

	billingAccount, err := j.repos.BillingAccounts.List(ctx)
	if err != nil {
		log.Printf("Error occurred while fetching customers: %v", err)
		return err
	}

	for _, ba := range billingAccount {
		for page := 0; ; page++ {
			resources, err := j.repos.Resources.ListByBillingAccountID(ctx, ba.AccountID, j.batchMaxSize, page*j.batchMaxSize)
			if err != nil {
				log.Printf("Error fetching resources for billing account %s: %v", ba.AccountID, err)
				return err
			}
			if len(resources) == 0 {
				break
			}

			if err := j.processBillingAccountResourcesBatch(ctx, startTime, ba, resources); err != nil {
				log.Printf("Error processing resources for billing account %s: %v", ba.AccountID, err)
				return err
			}
		}
	}

	return nil
}

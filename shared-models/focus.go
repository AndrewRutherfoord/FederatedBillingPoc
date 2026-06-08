package sharedmodels

import (
	"time"

	"github.com/shopspring/decimal"
)

type CapacityReservationStatus string

const (
	CapacityReservationStatusUsed   CapacityReservationStatus = "used"
	CapacityReservationStatusUnused CapacityReservationStatus = "unused"
)

type CommitmentDiscountCategory string

const (
	CommitmentDiscountCategorySpend CommitmentDiscountCategory = "spend"
	CommitmentDiscountCategoryUsage CommitmentDiscountCategory = "usage"
)

type CommitmentDiscountStatus string

const (
	CommitmentDiscountStatusUsed   CommitmentDiscountStatus = "used"
	CommitmentDiscountStatusUnused CommitmentDiscountStatus = "unused"
)

type ServiceCategory string

const (
	ServiceCategoryAIAndMachineLearning  ServiceCategory = "ai_and_machine_learning"
	ServiceCategoryAnalytics             ServiceCategory = "analytics"
	ServiceCategoryBusinessApplications  ServiceCategory = "business_applications"
	ServiceCategoryCompute               ServiceCategory = "compute"
	ServiceCategoryDatabases             ServiceCategory = "databases"
	ServiceCategoryDeveloperTools        ServiceCategory = "developer_tools"
	ServiceCategoryMulticloud            ServiceCategory = "multicloud"
	ServiceCategoryIdentity              ServiceCategory = "identity"
	ServiceCategoryIntegration           ServiceCategory = "integration"
	ServiceCategoryInternetOfThings      ServiceCategory = "internet_of_things"
	ServiceCategoryManagementGovernance  ServiceCategory = "management_and_governance"
	ServiceCategoryMedia                 ServiceCategory = "media"
	ServiceCategoryMigration             ServiceCategory = "migration"
	ServiceCategoryMobile                ServiceCategory = "mobile"
	ServiceCategoryNetworking            ServiceCategory = "networking"
	ServiceCategorySecurity              ServiceCategory = "security"
	ServiceCategoryStorage               ServiceCategory = "storage"
	ServiceCategoryWeb                   ServiceCategory = "web"
	ServiceCategoryOther                 ServiceCategory = "other"
)

type ServiceSubcategory string

const (
	// AI and Machine Learning
	ServiceSubcategoryAIPlatforms                ServiceSubcategory = "ai_platforms"
	ServiceSubcategoryBots                       ServiceSubcategory = "bots"
	ServiceSubcategoryGenerativeAI               ServiceSubcategory = "generative_ai"
	ServiceSubcategoryMachineLearning            ServiceSubcategory = "machine_learning"
	ServiceSubcategoryNaturalLanguageProcessing  ServiceSubcategory = "natural_language_processing"
	ServiceSubcategoryOtherAIAndML               ServiceSubcategory = "other_ai_and_machine_learning"

	// Analytics
	ServiceSubcategoryAnalyticsPlatforms ServiceSubcategory = "analytics_platforms"
	ServiceSubcategoryBusinessIntelligence ServiceSubcategory = "business_intelligence"
	ServiceSubcategoryDataProcessing     ServiceSubcategory = "data_processing"
	ServiceSubcategorySearch             ServiceSubcategory = "search"
	ServiceSubcategoryStreamingAnalytics ServiceSubcategory = "streaming_analytics"
	ServiceSubcategoryOtherAnalytics     ServiceSubcategory = "other_analytics"

	// Business Applications
	ServiceSubcategoryProductivityCollaboration ServiceSubcategory = "productivity_and_collaboration"
	ServiceSubcategoryOtherBusinessApplications ServiceSubcategory = "other_business_applications"

	// Compute
	ServiceSubcategoryContainers        ServiceSubcategory = "containers"
	ServiceSubcategoryEndUserComputing  ServiceSubcategory = "end_user_computing"
	ServiceSubcategoryQuantumCompute    ServiceSubcategory = "quantum_compute"
	ServiceSubcategoryServerlessCompute ServiceSubcategory = "serverless_compute"
	ServiceSubcategoryVirtualMachines   ServiceSubcategory = "virtual_machines"
	ServiceSubcategoryOtherCompute      ServiceSubcategory = "other_compute"

	// Databases
	ServiceSubcategoryCaching            ServiceSubcategory = "caching"
	ServiceSubcategoryDataWarehouses     ServiceSubcategory = "data_warehouses"
	ServiceSubcategoryLedgerDatabases    ServiceSubcategory = "ledger_databases"
	ServiceSubcategoryNoSQLDatabases     ServiceSubcategory = "nosql_databases"
	ServiceSubcategoryRelationalDatabases ServiceSubcategory = "relational_databases"
	ServiceSubcategoryTimeSeriesDatabases ServiceSubcategory = "time_series_databases"
	ServiceSubcategoryOtherDatabases     ServiceSubcategory = "other_databases"

	// Developer Tools
	ServiceSubcategoryDeveloperPlatforms          ServiceSubcategory = "developer_platforms"
	ServiceSubcategoryCICD                        ServiceSubcategory = "continuous_integration_and_deployment"
	ServiceSubcategoryDevelopmentEnvironments     ServiceSubcategory = "development_environments"
	ServiceSubcategorySourceCodeManagement        ServiceSubcategory = "source_code_management"
	ServiceSubcategoryQualityAssurance            ServiceSubcategory = "quality_assurance"
	ServiceSubcategoryOtherDeveloperTools         ServiceSubcategory = "other_developer_tools"

	// Identity
	ServiceSubcategoryIAM           ServiceSubcategory = "identity_and_access_management"
	ServiceSubcategoryOtherIdentity ServiceSubcategory = "other_identity"

	// Integration
	ServiceSubcategoryAPIManagement        ServiceSubcategory = "api_management"
	ServiceSubcategoryMessaging            ServiceSubcategory = "messaging"
	ServiceSubcategoryWorkflowOrchestration ServiceSubcategory = "workflow_orchestration"
	ServiceSubcategoryOtherIntegration     ServiceSubcategory = "other_integration"

	// Internet of Things
	ServiceSubcategoryIoTAnalytics        ServiceSubcategory = "iot_analytics"
	ServiceSubcategoryIoTPlatforms        ServiceSubcategory = "iot_platforms"
	ServiceSubcategoryOtherIoT            ServiceSubcategory = "other_internet_of_things"

	// Management and Governance
	ServiceSubcategoryArchitecture          ServiceSubcategory = "architecture"
	ServiceSubcategoryCompliance            ServiceSubcategory = "compliance"
	ServiceSubcategoryCostManagement        ServiceSubcategory = "cost_management"
	ServiceSubcategoryDataGovernance        ServiceSubcategory = "data_governance"
	ServiceSubcategoryDisasterRecovery      ServiceSubcategory = "disaster_recovery"
	ServiceSubcategoryEndpointManagement    ServiceSubcategory = "endpoint_management"
	ServiceSubcategoryObservability         ServiceSubcategory = "observability"
	ServiceSubcategorySupport               ServiceSubcategory = "support"
	ServiceSubcategoryOtherMgmtGovernance   ServiceSubcategory = "other_management_and_governance"

	// Media
	ServiceSubcategoryContentCreation ServiceSubcategory = "content_creation"
	ServiceSubcategoryGaming          ServiceSubcategory = "gaming"
	ServiceSubcategoryMediaStreaming  ServiceSubcategory = "media_streaming"
	ServiceSubcategoryMixedReality    ServiceSubcategory = "mixed_reality"
	ServiceSubcategoryOtherMedia      ServiceSubcategory = "other_media"

	// Migration
	ServiceSubcategoryDataMigration     ServiceSubcategory = "data_migration"
	ServiceSubcategoryResourceMigration ServiceSubcategory = "resource_migration"
	ServiceSubcategoryOtherMigration    ServiceSubcategory = "other_migration"

	// Mobile
	ServiceSubcategoryOtherMobile ServiceSubcategory = "other_mobile"

	// Multicloud
	ServiceSubcategoryMulticloudIntegration ServiceSubcategory = "multicloud_integration"
	ServiceSubcategoryOtherMulticloud       ServiceSubcategory = "other_multicloud"

	// Networking
	ServiceSubcategoryApplicationNetworking ServiceSubcategory = "application_networking"
	ServiceSubcategoryContentDelivery       ServiceSubcategory = "content_delivery"
	ServiceSubcategoryNetworkConnectivity   ServiceSubcategory = "network_connectivity"
	ServiceSubcategoryNetworkInfrastructure ServiceSubcategory = "network_infrastructure"
	ServiceSubcategoryNetworkRouting        ServiceSubcategory = "network_routing"
	ServiceSubcategoryNetworkSecurity       ServiceSubcategory = "network_security"
	ServiceSubcategoryOtherNetworking       ServiceSubcategory = "other_networking"

	// Security
	ServiceSubcategorySecretManagement         ServiceSubcategory = "secret_management"
	ServiceSubcategorySecurityPostureMgmt      ServiceSubcategory = "security_posture_management"
	ServiceSubcategoryThreatDetectionResponse  ServiceSubcategory = "threat_detection_and_response"
	ServiceSubcategoryOtherSecurity            ServiceSubcategory = "other_security"

	// Storage
	ServiceSubcategoryBackupStorage   ServiceSubcategory = "backup_storage"
	ServiceSubcategoryBlockStorage    ServiceSubcategory = "block_storage"
	ServiceSubcategoryFileStorage     ServiceSubcategory = "file_storage"
	ServiceSubcategoryObjectStorage   ServiceSubcategory = "object_storage"
	ServiceSubcategoryStoragePlatforms ServiceSubcategory = "storage_platforms"
	ServiceSubcategoryOtherStorage    ServiceSubcategory = "other_storage"

	// Web
	ServiceSubcategoryApplicationPlatforms ServiceSubcategory = "application_platforms"
	ServiceSubcategoryOtherWeb             ServiceSubcategory = "other_web"

	// Other
	ServiceSubcategoryOtherOther ServiceSubcategory = "other_other"
)

// FocusLineItem represents a charge according to the FinOps FOCUS 1.3 specification.
// Reference: https://focus.finops.org/focus-specification/1-3/
type FocusLineItem struct {
	// Required Fields
	BilledCost          decimal.Decimal `json:"billed_cost"`
	BillingAccountID    string          `json:"billing_account_id"`
	BillingAccountType  string          `json:"billing_account_type"`
	BillingCurrency     string          `json:"billing_currency"`
	BillingPeriodEnd    time.Time       `json:"billing_period_end"`
	BillingPeriodStart  time.Time       `json:"billing_period_start"`
	ChargeCategory      string          `json:"charge_category"`
	ChargeFrequency     string          `json:"charge_frequency"`
	ChargePeriodEnd     time.Time       `json:"charge_period_end"`
	ChargePeriodStart   time.Time       `json:"charge_period_start"`
	ContractedCost      decimal.Decimal `json:"contracted_cost"`
	EffectiveCost       decimal.Decimal `json:"effective_cost"`
	HostProviderName    string          `json:"host_provider_name"`
	InvoiceIssuerName   string          `json:"invoice_issuer_name"`
	ListCost            decimal.Decimal `json:"list_cost"`
	ServiceCategory     ServiceCategory `json:"service_category"`
	ServiceName         string          `json:"service_name"`
	ServiceProviderName string          `json:"service_provider_name"`
	ServiceSubcategory  ServiceSubcategory `json:"service_subcategory"`

	// Optional Fields - Allocation
	AllocatedMethodDetails *string `json:"allocated_method_details,omitempty"`
	AllocatedMethodID      *string `json:"allocated_method_id,omitempty"`
	AllocatedResourceID    *string `json:"allocated_resource_id,omitempty"`
	AllocatedResourceName  *string `json:"allocated_resource_name,omitempty"`
	AllocatedTags          *string `json:"allocated_tags,omitempty"`

	// Optional Fields - Location
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	RegionID         *string `json:"region_id,omitempty"`
	RegionName       *string `json:"region_name,omitempty"`

	// Optional Fields - Billing Account
	BillingAccountName *string `json:"billing_account_name,omitempty"`

	// Optional Fields - Capacity Reservation
	CapacityReservationID     *string                    `json:"capacity_reservation_id,omitempty"`
	CapacityReservationStatus *CapacityReservationStatus `json:"capacity_reservation_status,omitempty"`

	// Optional Fields - Charge Details
	ChargeClass       *string `json:"charge_class,omitempty"`
	ChargeDescription *string `json:"charge_description,omitempty"`

	// Optional Fields - Commitment Discount
	CommitmentDiscountCategory *CommitmentDiscountCategory `json:"commitment_discount_category,omitempty"`
	CommitmentDiscountID       *string                     `json:"commitment_discount_id,omitempty"`
	CommitmentDiscountName     *string                     `json:"commitment_discount_name,omitempty"`
	CommitmentDiscountQuantity *decimal.Decimal            `json:"commitment_discount_quantity,omitempty"`
	CommitmentDiscountStatus   *CommitmentDiscountStatus   `json:"commitment_discount_status,omitempty"`
	CommitmentDiscountType     *string                     `json:"commitment_discount_type,omitempty"`
	CommitmentDiscountUnit     *string                     `json:"commitment_discount_unit,omitempty"`

	// Optional Fields - Consumption
	ConsumedQuantity *decimal.Decimal `json:"consumed_quantity,omitempty"`
	ConsumedUnit     *string          `json:"consumed_unit,omitempty"`

	// Optional Fields - Contract
	ContractApplied      *string          `json:"contract_applied,omitempty"`
	ContractedUnitPrice  *decimal.Decimal `json:"contracted_unit_price,omitempty"`

	// Optional Fields - Invoice
	InvoiceID *string `json:"invoice_id,omitempty"`

	// Optional Fields - Pricing
	ListUnitPrice                       *decimal.Decimal `json:"list_unit_price,omitempty"`
	PricingCategory                     *string          `json:"pricing_category,omitempty"`
	PricingCurrency                     *string          `json:"pricing_currency,omitempty"`
	PricingCurrencyContractedUnitPrice  *decimal.Decimal `json:"pricing_currency_contracted_unit_price,omitempty"`
	PricingCurrencyEffectiveCost        *decimal.Decimal `json:"pricing_currency_effective_cost,omitempty"`
	PricingCurrencyListUnitPrice        *decimal.Decimal `json:"pricing_currency_list_unit_price,omitempty"`
	PricingQuantity                     *decimal.Decimal `json:"pricing_quantity,omitempty"`
	PricingUnit                         *string          `json:"pricing_unit,omitempty"`

	// Optional Fields - Resource
	ResourceID   *string `json:"resource_id,omitempty"`
	ResourceName *string `json:"resource_name,omitempty"`
	ResourceType *string `json:"resource_type,omitempty"`

	// Optional Fields - SKU
	SkuID           *string `json:"sku_id,omitempty"`
	SkuMeter        *string `json:"sku_meter,omitempty"`
	SkuPriceDetails *string `json:"sku_price_details,omitempty"`
	SkuPriceID      *string `json:"sku_price_id,omitempty"`

	// Optional Fields - Sub Account
	SubAccountID   *string `json:"sub_account_id,omitempty"`
	SubAccountName *string `json:"sub_account_name,omitempty"`
	SubAccountType *string `json:"sub_account_type,omitempty"`

	// Optional Fields - Tags
	Tags *string `json:"tags,omitempty"`
}

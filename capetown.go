// Package capetown provides named layer locators and pre-built QueryParams for
// the City of Cape Town Open Data Portal.
//
// In 2025 the City retired the single monolithic
// Theme_Based/Open_Data_Service and split its layers across a set of themed
// feature services named ODP_SPLIT_1 … ODP_SPLIT_13, renumbering layer IDs per
// service from zero. A dataset is therefore identified by BOTH a service name
// and a layer ID, so the constructors here return a [Query] carrying both.
//
// The pre-built queries deliberately do not pin an output field list. The
// upstream layer schemas drift and use non-obvious column names, so each query
// returns the layer's full field set; callers select fields explicitly when
// they want a smaller payload. Filters (suburb, ordering) reference field names
// verified against the live service.
//
// Service names and layer IDs are best-effort and may drift as the portal is
// republished. Confirm a suspect locator against the live service with a
// [github.com/richardwooding/go-arcgis.Client] ServiceInfo/LayerInfo before
// relying on it.
package capetown

import (
	"fmt"
	"strings"

	arcgis "github.com/richardwooding/go-arcgis"
)

// The City publishes the ODP_SPLIT services from two ArcGIS servers. The open
// data hub links to esapqa, which also carries ODP_SPLIT_9 as a feature service
// and ODP_SPLIT_13; citymaps has neither. Record counts differ on some layers.
const (
	FolderESAPQA   = "https://esapqa.capetown.gov.za/agsext/rest/services/Theme_Based"
	FolderCityMaps = "https://citymaps.capetown.gov.za/agsext/rest/services/Theme_Based"
)

// BaseFolder is the default ArcGIS REST folder for the ODP_SPLIT services.
const BaseFolder = FolderESAPQA

// ServiceURL returns the FeatureServer endpoint for a service on BaseFolder.
func ServiceURL(service string) string {
	return ServiceURLIn(BaseFolder, service)
}

// ServiceURLIn returns the FeatureServer endpoint for a service under folder.
// ArcGIS Online datasets ignore folder and return their fallback URL.
func ServiceURLIn(folder, service string) string {
	if h, ok := hubDatasets[service]; ok {
		return h.fallbackURL
	}
	return strings.TrimRight(folder, "/") + "/" + service + "/FeatureServer"
}

// Datasets the City hosts on ArcGIS Online rather than the ODP_SPLIT services.
// Their service names embed dates and go stale, so each is pinned to its hub
// item ID; resolve the current URL with arcgis.ItemURL and fall back to
// ServiceURL.
const (
	ServiceServiceRequests     = "SERVICE_REQUESTS"
	ServiceBuildingPlans       = "BUILDING_PLANS"
	ServicePopulationEstimates = "POPULATION_ESTIMATES"
	ServiceTaxiFares           = "TAXI_FARES"
)

// Layer IDs within the ArcGIS Online datasets; all are non-spatial tables.
const (
	LayerServiceRequests     = 0
	LayerBuildingPlans       = 0
	LayerPopulationEstimates = 0
	LayerTaxiFares           = 0
)

type hubDataset struct {
	itemID      string
	fallbackURL string
}

const agolFolder = "https://services6.arcgis.com/nyYfO9SxHU2ChQd9/arcgis/rest/services"

var hubDatasets = map[string]hubDataset{
	// The City froze the original service requests item
	// (90711027d33940b5a06e96ad8a8f7ede) on 6 October 2026 and publishes the
	// live table under a new item with a reworked schema (CreatedOnDT rather
	// than Created_On_Date, plus Suburb_Clean, Longitude and Latitude).
	ServiceServiceRequests: {
		itemID:      "da3b6d46845e4671bde500e8de7e1d3e",
		fallbackURL: agolFolder + "/schema_jan2023_iso/FeatureServer",
	},
	ServiceBuildingPlans: {
		itemID:      "a9a32e2ecc074b318ee0ad136fb7e2af",
		fallbackURL: agolFolder + "/Building_Plan_Approvals_2014_to_2025/FeatureServer",
	},
	ServicePopulationEstimates: {
		itemID:      "61bd1a6c374d44a093862784b79f6c23",
		fallbackURL: agolFolder + "/Population_Estimates_2026_(Suburb_Level)/FeatureServer",
	},
	ServiceTaxiFares: {
		itemID:      "d160df6a7d7d4006a3eacf265585a27a",
		fallbackURL: agolFolder + "/Minibus_taxi_fares_surveyed_between_2023_2025_Rev2_06112026/FeatureServer",
	},
}

// HubItemID returns the ArcGIS Online item ID for a hub-hosted service key.
func HubItemID(service string) (string, bool) {
	h, ok := hubDatasets[service]
	return h.itemID, ok
}

// HubServices lists the service keys hosted on ArcGIS Online.
func HubServices() []string {
	return []string{ServiceServiceRequests, ServiceBuildingPlans, ServicePopulationEstimates, ServiceTaxiFares}
}

// Split feature-service names hosting the well-known datasets, validated
// against the live portal.
const (
	ServicePublicLighting    = "ODP_SPLIT_1"
	ServiceHeritageInventory = "ODP_SPLIT_3"
	ServiceLandParcels       = "ODP_SPLIT_4"
	ServiceWards             = "ODP_SPLIT_5"
	ServiceOfficialSuburbs   = "ODP_SPLIT_5"
	ServiceTaxiRoutes        = "ODP_SPLIT_6"
	ServiceLoadShedding      = "ODP_SPLIT_7"
	ServiceWaterQuality      = "ODP_SPLIT_12"
)

// Layer IDs for well-known CCT datasets, numbered within their split service.
const (
	LayerPublicLighting     = 1  // ODP_SPLIT_1 "Electricity Public Lighting"
	LayerHeritageInventory  = 2  // ODP_SPLIT_3 "Heritage Inventory"
	LayerLandParcels        = 0  // ODP_SPLIT_4 "Land Parcels"
	LayerWards              = 6  // ODP_SPLIT_5 "Ward"
	LayerOfficialSuburbs    = 3  // ODP_SPLIT_5 "Official Suburb"
	LayerTaxiRoutes         = 11 // ODP_SPLIT_6 "Taxi Routes"
	LayerLoadSheddingBlocks = 13 // ODP_SPLIT_7 "Loadshedding Blocks"
	LayerWaterQuality       = 12 // ODP_SPLIT_12 "Inland Water Quality Results (Raw)" (a table)
)

// Query locates a dataset: the split service that hosts it plus the ArcGIS
// query parameters (layer ID, filters, ordering) to run against that service.
type Query struct {
	// Service is the ODP_SPLIT_* feature service name (feed it to ServiceURL).
	Service string
	// Params are the ArcGIS query parameters, including the in-service LayerID.
	Params arcgis.QueryParams
}

// Services returns the canonical list of split feature services that make up
// the Open Data portal. Callers that want to enumerate everything (e.g. an
// aggregated service listing) iterate this set. It is best-effort: individual
// services may be temporarily unavailable as the portal is republished.
func Services() []string {
	return []string{
		"ODP_SPLIT_1", "ODP_SPLIT_2", "ODP_SPLIT_3", "ODP_SPLIT_4",
		"ODP_SPLIT_5", "ODP_SPLIT_6", "ODP_SPLIT_7", "ODP_SPLIT_8",
		"ODP_SPLIT_9", "ODP_SPLIT_10", "ODP_SPLIT_11", "ODP_SPLIT_12",
		"ODP_SPLIT_13",
	}
}

// fieldOfficialSuburbName is the official-suburb-name column shared by the
// land parcels and official suburb layers. Other CCT layers name their suburb
// column differently, so do not assume it elsewhere.
const fieldOfficialSuburbName = "OFC_SBRB_NAME"

// equalsFold builds a case-insensitive equality clause. Several CCT layers
// store names in upper case while callers pass mixed case, so a naive equality
// would silently return nothing. Single quotes in value are escaped.
func equalsFold(field, value string) string {
	return fmt.Sprintf("UPPER(%s) = UPPER('%s')", field, quoteSQL(value))
}

// quoteSQL escapes a string literal for an ArcGIS where clause.
func quoteSQL(s string) string { return strings.ReplaceAll(s, "'", "''") }

// --- Load Shedding ---

// LoadSheddingBlocks returns all load shedding block polygons. The layer carries
// only block geometry and a BlockID; it has no stage or suburb attribute.
func LoadSheddingBlocks() Query {
	return Query{Service: ServiceLoadShedding, Params: arcgis.QueryParams{LayerID: LayerLoadSheddingBlocks}}
}

// --- Wards ---

// Wards returns all municipal ward boundaries.
func Wards() Query {
	return Query{Service: ServiceWards, Params: arcgis.QueryParams{LayerID: LayerWards}}
}

// --- Official Suburbs ---

// OfficialSuburbs returns the official suburb boundary polygons, ordered by
// name. OFC_SBRB_NAME is upper case; it is the key that joins suburb-level
// tables such as PopulationEstimates to a shape.
func OfficialSuburbs() Query {
	return Query{
		Service: ServiceOfficialSuburbs,
		Params: arcgis.QueryParams{
			LayerID:       LayerOfficialSuburbs,
			OrderByFields: []string{fieldOfficialSuburbName},
		},
	}
}

// OfficialSuburbByName returns the official suburb polygon(s) matching a name,
// case-insensitively.
func OfficialSuburbByName(suburb string) Query {
	q := OfficialSuburbs()
	q.Params.Where = equalsFold(fieldOfficialSuburbName, suburb)
	return q
}

// --- Land Parcels ---

// LandParcels returns cadastral land parcel (erf) polygons.
func LandParcels() Query {
	return Query{Service: ServiceLandParcels, Params: arcgis.QueryParams{LayerID: LayerLandParcels}}
}

// LandParcelsBySuburb filters land parcels by official suburb name. The match
// is case-insensitive: the upstream column stores suburb names in upper case
// (e.g. "NEWLANDS"), so a naive equality on mixed-case input would silently
// return nothing.
func LandParcelsBySuburb(suburb string) Query {
	q := LandParcels()
	q.Params.Where = equalsFold(fieldOfficialSuburbName, suburb)
	return q
}

// --- Transport ---

// TaxiRoutes returns all registered minibus taxi routes.
func TaxiRoutes() Query {
	return Query{Service: ServiceTaxiRoutes, Params: arcgis.QueryParams{LayerID: LayerTaxiRoutes}}
}

// TaxiFares returns minibus taxi fares surveyed at ranks, newest survey first.
// Each row is an origin-destination rank pair with the average fare in rands
// (Fare), the survey date and, for each end, a rank reference (e.g. MBT019), a
// name and WGS84 longitude/latitude held as plain attribute columns; the table
// has no geometry.
func TaxiFares() Query {
	return Query{
		Service: ServiceTaxiFares,
		Params: arcgis.QueryParams{
			LayerID:       LayerTaxiFares,
			OrderByFields: []string{"Survey_Date DESC"},
		},
	}
}

// TaxiFaresFromRank filters surveyed fares to those departing a rank, by its
// reference code (Origin_Rank_Ref_No, e.g. "MBT019"). The codes are the stable
// key; rank names carry inconsistent spacing.
func TaxiFaresFromRank(originRankRef string) Query {
	q := TaxiFares()
	q.Params.Where = fmt.Sprintf("Origin_Rank_Ref_No = '%s'", quoteSQL(originRankRef))
	return q
}

// --- Electricity ---

// PublicLighting returns public street-lighting assets.
func PublicLighting() Query {
	return Query{Service: ServicePublicLighting, Params: arcgis.QueryParams{LayerID: LayerPublicLighting}}
}

// --- Heritage ---

// HeritageInventory returns heritage inventory sites and features.
func HeritageInventory() Query {
	return Query{Service: ServiceHeritageInventory, Params: arcgis.QueryParams{LayerID: LayerHeritageInventory}}
}

// --- Water ---

// WaterQualityResults returns inland water quality measurements, most recent
// first. It targets the raw results table (sample point, date, parameter,
// value); the table has no geometry.
func WaterQualityResults() Query {
	return Query{
		Service: ServiceWaterQuality,
		Params: arcgis.QueryParams{
			LayerID:       LayerWaterQuality,
			OrderByFields: []string{"SMPL_DATE DESC"},
		},
	}
}

// --- Demographics ---

// PopulationEstimates returns the City's model-based population and household
// estimates per official suburb, ordered by suburb name. They are an annual
// planning estimate derived from the 2022 census, not an official count, and
// the values are fractional. Suburb_Name is mixed case; match it with
// PopulationEstimateForSuburb or join it case-insensitively to OFC_SBRB_NAME
// on OfficialSuburbs.
func PopulationEstimates() Query {
	return Query{
		Service: ServicePopulationEstimates,
		Params: arcgis.QueryParams{
			LayerID:       LayerPopulationEstimates,
			OrderByFields: []string{"Suburb_Name"},
		},
	}
}

// PopulationEstimateForSuburb returns the estimate row for one suburb,
// matched case-insensitively.
func PopulationEstimateForSuburb(suburb string) Query {
	q := PopulationEstimates()
	q.Params.Where = equalsFold("Suburb_Name", suburb)
	return q
}

// --- ArcGIS Online datasets ---

// ServiceRequests returns citizen service requests, most recently loaded first.
// Sorting the four-and-a-half-million-row table by CreatedOnDT takes ArcGIS
// Online tens of seconds and times out when paging, whereas the object ID is
// indexed and its highest values are the newest requests.
func ServiceRequests() Query {
	return Query{
		Service: ServiceServiceRequests,
		Params: arcgis.QueryParams{
			LayerID:       LayerServiceRequests,
			OrderByFields: []string{"ObjectId DESC"},
		},
	}
}

// BuildingPlanApprovals returns building plan applications, newest submission first.
func BuildingPlanApprovals() Query {
	return Query{
		Service: ServiceBuildingPlans,
		Params: arcgis.QueryParams{
			LayerID:       LayerBuildingPlans,
			OrderByFields: []string{"Submission_Date DESC"},
		},
	}
}

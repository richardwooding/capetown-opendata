// Package capetown provides named layer locators and pre-built QueryParams for
// the City of Cape Town Open Data Portal.
//
// In 2025 the City retired the single monolithic
// Theme_Based/Open_Data_Service and split its layers across a set of themed
// feature services named ODP_SPLIT_1 … ODP_SPLIT_12, renumbering layer IDs per
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

// BaseFolder is the ArcGIS REST folder that hosts the City of Cape Town Open
// Data feature services (the ODP_SPLIT_* services live directly under it).
const BaseFolder = "https://citymaps.capetown.gov.za/agsext/rest/services/Theme_Based"

// ServiceURL returns the FeatureServer endpoint for a named split service, e.g.
// ServiceURL("ODP_SPLIT_5") →
// https://citymaps.capetown.gov.za/agsext/rest/services/Theme_Based/ODP_SPLIT_5/FeatureServer.
func ServiceURL(service string) string {
	return BaseFolder + "/" + service + "/FeatureServer"
}

// Split feature-service names hosting the well-known datasets, validated
// against the live portal.
const (
	ServicePublicLighting    = "ODP_SPLIT_1"
	ServiceHeritageInventory = "ODP_SPLIT_3"
	ServiceLandParcels       = "ODP_SPLIT_4"
	ServiceWards             = "ODP_SPLIT_5"
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
	}
}

// fieldLandParcelSuburb is the official-suburb-name column on the land parcels
// layer. The suburb column name is not consistent across CCT layers, so it is
// scoped to the dataset that uses it.
const fieldLandParcelSuburb = "OFC_SBRB_NAME"

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
	q.Params.Where = fmt.Sprintf("UPPER(%s) = UPPER('%s')", fieldLandParcelSuburb, strings.ReplaceAll(suburb, "'", "''"))
	return q
}

// --- Transport ---

// TaxiRoutes returns all registered minibus taxi routes.
func TaxiRoutes() Query {
	return Query{Service: ServiceTaxiRoutes, Params: arcgis.QueryParams{LayerID: LayerTaxiRoutes}}
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

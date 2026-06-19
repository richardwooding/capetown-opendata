// Package capetown provides named layer IDs and pre-built QueryParams
// for the City of Cape Town Open Data Portal.
//
// Base URL: https://citymaps.capetown.gov.za/agsext/rest/services/Theme_Based/Open_Data_Service/FeatureServer
//
// The pre-built queries deliberately do not pin an output field list. The
// upstream layer schemas drift and use non-obvious column names, so each query
// returns the layer's full field set; callers select fields explicitly when
// they want a smaller payload. Filters (suburb, ordering) reference field names
// verified against the live service.
//
// Layer IDs are best-effort and may drift as the service is republished. The
// integration test (run with -tags=integration) validates every layer ID and
// filter field against the live service. Use a
// [github.com/richardwooding/go-arcgis.Client] with ServiceInfo to confirm them
// at runtime.
package capetown

import (
	"fmt"
	"strings"

	arcgis "github.com/richardwooding/go-arcgis"
)

// BaseURL is the City of Cape Town Open Data Feature Service endpoint.
const BaseURL = "https://citymaps.capetown.gov.za/agsext/rest/services/Theme_Based/Open_Data_Service/FeatureServer"

// Layer IDs for well-known CCT datasets, validated against the live service.
const (
	LayerLoadSheddingBlocks = 138
	LayerWards              = 78
	LayerLandParcels        = 56
	LayerTaxiRoutes         = 97
	LayerPublicLighting     = 3
	LayerWaterQuality       = 229
	LayerHeritageInventory  = 49
)

// fieldLandParcelSuburb is the official-suburb-name column on the land parcels
// layer. The suburb column name is not consistent across CCT layers, so it is
// scoped to the dataset that uses it.
const fieldLandParcelSuburb = "OFC_SBRB_NAME"

// --- Load Shedding ---

// LoadSheddingBlocks returns all load shedding block polygons. The layer carries
// only block geometry and a BlockID; it has no stage or suburb attribute.
func LoadSheddingBlocks() arcgis.QueryParams {
	return arcgis.QueryParams{LayerID: LayerLoadSheddingBlocks}
}

// --- Wards ---

// Wards returns all municipal ward boundaries.
func Wards() arcgis.QueryParams {
	return arcgis.QueryParams{LayerID: LayerWards}
}

// --- Land Parcels ---

// LandParcels returns cadastral land parcel (erf) polygons.
func LandParcels() arcgis.QueryParams {
	return arcgis.QueryParams{LayerID: LayerLandParcels}
}

// LandParcelsBySuburb filters land parcels by official suburb name.
func LandParcelsBySuburb(suburb string) arcgis.QueryParams {
	p := LandParcels()
	p.Where = fmt.Sprintf("%s = '%s'", fieldLandParcelSuburb, strings.ReplaceAll(suburb, "'", "''"))
	return p
}

// --- Transport ---

// TaxiRoutes returns all registered minibus taxi routes.
func TaxiRoutes() arcgis.QueryParams {
	return arcgis.QueryParams{LayerID: LayerTaxiRoutes}
}

// --- Water ---

// WaterQualityResults returns inland water quality measurements, most recent
// first. It targets the raw results table (sample point, date, parameter,
// value).
func WaterQualityResults() arcgis.QueryParams {
	return arcgis.QueryParams{
		LayerID:       LayerWaterQuality,
		OrderByFields: []string{"SMPL_DATE DESC"},
	}
}

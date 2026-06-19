//go:build integration

// Package capetown integration tests validate the pre-built queries against the
// live City of Cape Town Feature Service. They hit the network and are excluded
// from the default build; run them with:
//
//	go test -tags=integration ./...
//
// Their job is to catch upstream drift: a layer ID that has moved, a filter or
// order-by column that has been renamed, or a query the service rejects.
package capetown_test

import (
	"context"
	"testing"
	"time"

	arcgis "github.com/richardwooding/go-arcgis"

	capetown "github.com/richardwooding/capetown-opendata"
)

// liveTimeout is generous: the live municipal service has variable latency, and
// these tests exist to catch schema drift, not to benchmark response times.
const liveTimeout = 60 * time.Second

func liveClient() *arcgis.Client {
	return arcgis.NewClient(capetown.BaseURL, arcgis.WithTimeout(liveTimeout))
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), liveTimeout)
	t.Cleanup(cancel)
	return c
}

// TestLiveLayerIDsExist asserts every named layer ID is still published by the
// service (as either a layer or a table).
func TestLiveLayerIDsExist(t *testing.T) {
	info, err := liveClient().ServiceInfo(ctx(t))
	if err != nil {
		t.Fatalf("ServiceInfo: %v", err)
	}
	present := map[int]bool{}
	for _, l := range info.Layers {
		present[l.ID] = true
	}
	for _, tbl := range info.Tables {
		present[tbl.ID] = true
	}
	for name, id := range map[string]int{
		"LayerLoadSheddingBlocks": capetown.LayerLoadSheddingBlocks,
		"LayerWards":              capetown.LayerWards,
		"LayerLandParcels":        capetown.LayerLandParcels,
		"LayerTaxiRoutes":         capetown.LayerTaxiRoutes,
		"LayerPublicLighting":     capetown.LayerPublicLighting,
		"LayerWaterQuality":       capetown.LayerWaterQuality,
		"LayerHeritageInventory":  capetown.LayerHeritageInventory,
	} {
		if !present[id] {
			t.Errorf("%s = %d is no longer published by the service", name, id)
		}
	}
}

// TestLiveNamedQueriesSucceed runs every pre-built query against the live
// service. A drifted layer ID, bad order-by field, or otherwise malformed query
// surfaces here as a non-nil error.
func TestLiveNamedQueriesSucceed(t *testing.T) {
	c := liveClient()
	cases := map[string]arcgis.QueryParams{
		"LoadSheddingBlocks":  capetown.LoadSheddingBlocks(),
		"Wards":               capetown.Wards(),
		"LandParcels":         capetown.LandParcels(),
		"LandParcelsBySuburb": capetown.LandParcelsBySuburb("Newlands"),
		"TaxiRoutes":          capetown.TaxiRoutes(),
		"WaterQualityResults": capetown.WaterQualityResults(),
	}
	for name, p := range cases {
		t.Run(name, func(t *testing.T) {
			p.PageSize = 1
			if _, err := c.Query(ctx(t), p); err != nil {
				t.Errorf("%s query failed against live service: %v", name, err)
			}
		})
	}
}

// TestLiveFilterFieldsExist asserts the columns referenced by filters and
// ordering still exist on their layers.
func TestLiveFilterFieldsExist(t *testing.T) {
	c := liveClient()
	cases := []struct {
		name    string
		layerID int
		field   string
	}{
		{"land parcel suburb", capetown.LayerLandParcels, "OFC_SBRB_NAME"},
		{"water quality sample date", capetown.LayerWaterQuality, "SMPL_DATE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := c.LayerInfo(ctx(t), tc.layerID)
			if err != nil {
				t.Fatalf("LayerInfo(%d): %v", tc.layerID, err)
			}
			for _, f := range info.Fields {
				if f.Name == tc.field {
					return
				}
			}
			t.Errorf("field %q not found on layer %d (%s)", tc.field, tc.layerID, info.Name)
		})
	}
}

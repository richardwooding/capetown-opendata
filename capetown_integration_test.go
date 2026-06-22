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
	"errors"
	"testing"
	"time"

	arcgis "github.com/richardwooding/go-arcgis"

	capetown "github.com/richardwooding/capetown-opendata"
)

// liveTimeout is generous: the live municipal service has variable latency, and
// these tests exist to catch schema drift, not to benchmark response times.
const liveTimeout = 60 * time.Second

// liveAttempts is how many times a live call is retried when it times out. The
// municipal service has intermittent multi-second latency spikes; without
// retries a single slow response turns the daily drift check red even though
// nothing has actually drifted. Genuine drift (a moved layer ID, a renamed
// field, an HTTP 4xx) surfaces as a deterministic error rather than a timeout,
// so it is not retried — see retry.
const liveAttempts = 3

func liveClient() *arcgis.Client {
	return arcgis.NewClient(capetown.BaseURL, arcgis.WithTimeout(liveTimeout))
}

// retry runs fn with a fresh per-attempt timeout context until it succeeds,
// returns a non-timeout error, or attempts are exhausted, returning the final
// result and error. Only context.DeadlineExceeded is treated as transient and
// retried; every other error (including the deterministic ones that signal real
// drift) returns immediately so the check still fails fast and loudly.
func retry[T any](t *testing.T, label string, fn func(context.Context) (T, error)) (T, error) {
	t.Helper()
	var (
		out T
		err error
	)
	for attempt := 1; attempt <= liveAttempts; attempt++ {
		c, cancel := context.WithTimeout(context.Background(), liveTimeout)
		out, err = fn(c)
		cancel()
		if err == nil || !errors.Is(err, context.DeadlineExceeded) {
			return out, err
		}
		t.Logf("%s: attempt %d/%d timed out, retrying: %v", label, attempt, liveAttempts, err)
	}
	return out, err
}

// TestLiveLayerIDsExist asserts every named layer ID is still published by the
// service (as either a layer or a table).
func TestLiveLayerIDsExist(t *testing.T) {
	info, err := retry(t, "ServiceInfo", liveClient().ServiceInfo)
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
			_, err := retry(t, name, func(c2 context.Context) (*arcgis.FeatureSet, error) {
				return c.Query(c2, p)
			})
			if err != nil {
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
			info, err := retry(t, tc.name, func(c2 context.Context) (*arcgis.LayerInfo, error) {
				return c.LayerInfo(c2, tc.layerID)
			})
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

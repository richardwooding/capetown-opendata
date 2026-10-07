//go:build integration

// Package capetown integration tests validate the pre-built queries against the
// live City of Cape Town feature services. They hit the network and are
// excluded from the default build; run them with:
//
//	go test -tags=integration ./...
//
// Their job is to catch upstream drift: a dataset that has moved to a different
// split service, a layer ID that has moved, a filter or order-by column that
// has been renamed, or a query the service rejects.
package capetown_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	arcgis "github.com/richardwooding/go-arcgis"

	capetown "github.com/richardwooding/capetown-opendata"
)

// liveTimeout is generous: the live municipal service has variable latency, and
// these tests exist to catch schema drift, not to benchmark response times.
const liveTimeout = 60 * time.Second

// liveAttempts is how many times a live call is retried on a transient failure.
// The municipal service has intermittent multi-second latency spikes and
// periodically returns HTTP 5xx while a MapServer restarts; without retries a
// single slow or server-side blip turns the daily drift check red even though
// nothing has actually drifted. Genuine drift (a moved layer ID, a renamed
// field, an HTTP 4xx) surfaces as a deterministic error rather than a transient
// one, so it is not retried — see retry and isTransient.
const liveAttempts = 3

// datasets is every pre-built dataset locator, keyed by a readable name.
func datasets() map[string]capetown.Query {
	return map[string]capetown.Query{
		"LoadSheddingBlocks":  capetown.LoadSheddingBlocks(),
		"Wards":               capetown.Wards(),
		"LandParcels":         capetown.LandParcels(),
		"LandParcelsBySuburb": capetown.LandParcelsBySuburb("Newlands"),
		"TaxiRoutes":          capetown.TaxiRoutes(),
		"PublicLighting":      capetown.PublicLighting(),
		"WaterQualityResults": capetown.WaterQualityResults(),
		"HeritageInventory":   capetown.HeritageInventory(),
		"ServiceRequests":     capetown.ServiceRequests(),
		"BuildingPlans":       capetown.BuildingPlanApprovals(),
	}
}

// liveClient targets BaseFolder unless CAPETOWN_BASE_FOLDER names another
// server folder, so the alternate City host can be checked by hand.
func liveClient(service string) *arcgis.Client {
	folder := capetown.BaseFolder
	if f := os.Getenv("CAPETOWN_BASE_FOLDER"); f != "" {
		folder = f
	}
	return arcgis.NewClient(capetown.ServiceURLIn(folder, service), arcgis.WithTimeout(liveTimeout))
}

// isTransient reports whether err is a temporary upstream condition that a
// retry might clear, as opposed to a deterministic error that signals real
// drift. Two cases qualify: a per-attempt timeout (latency spike) and an ArcGIS
// 5xx (the service is down or a MapServer is still starting). A 4xx — an
// invalid query, a missing layer — is deterministic and returns immediately so
// the check still fails fast and loudly.
func isTransient(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var apiErr *arcgis.APIError
	if errors.As(err, &apiErr) && apiErr.Code >= 500 && apiErr.Code < 600 {
		return true
	}
	return false
}

// retry runs fn with a fresh per-attempt timeout context until it succeeds,
// returns a non-transient error, or attempts are exhausted, returning the final
// result and error. Only transient errors (see isTransient) are retried; every
// other error returns immediately.
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
		if err == nil || !isTransient(err) {
			return out, err
		}
		t.Logf("%s: attempt %d/%d failed transiently, retrying: %v", label, attempt, liveAttempts, err)
	}
	return out, err
}

// TestLiveLayerIDsExist asserts every dataset's layer ID is still published by
// the split service it claims to live on (as either a layer or a table).
func TestLiveLayerIDsExist(t *testing.T) {
	for name, q := range datasets() {
		t.Run(name, func(t *testing.T) {
			c := liveClient(q.Service)
			info, err := retry(t, name, c.ServiceInfo)
			if err != nil {
				t.Fatalf("ServiceInfo(%s): %v", q.Service, err)
			}
			present := map[int]bool{}
			for _, l := range info.Layers {
				present[l.ID] = true
			}
			for _, tbl := range info.Tables {
				present[tbl.ID] = true
			}
			if !present[q.Params.LayerID] {
				t.Errorf("%s: layer %d is no longer published by %s", name, q.Params.LayerID, q.Service)
			}
		})
	}
}

// TestLiveNamedQueriesSucceed runs every pre-built query against its live split
// service. A drifted layer ID, bad order-by field, or otherwise malformed query
// surfaces here as a non-nil error.
func TestLiveNamedQueriesSucceed(t *testing.T) {
	for name, q := range datasets() {
		t.Run(name, func(t *testing.T) {
			c := liveClient(q.Service)
			p := q.Params
			p.PageSize = 1
			_, err := retry(t, name, func(c2 context.Context) (*arcgis.FeatureSet, error) {
				return c.Query(c2, p)
			})
			if err != nil {
				t.Errorf("%s query failed against live service %s: %v", name, q.Service, err)
			}
		})
	}
}

// TestLiveFilterFieldsExist asserts the columns referenced by filters and
// ordering still exist on their layers.
func TestLiveFilterFieldsExist(t *testing.T) {
	cases := []struct {
		name    string
		service string
		layerID int
		field   string
	}{
		{"land parcel suburb", capetown.ServiceLandParcels, capetown.LayerLandParcels, "OFC_SBRB_NAME"},
		{"water quality sample date", capetown.ServiceWaterQuality, capetown.LayerWaterQuality, "SMPL_DATE"},
		{"service request ward", capetown.ServiceServiceRequests, capetown.LayerServiceRequests, "Ward"},
		{"service request complaint type", capetown.ServiceServiceRequests, capetown.LayerServiceRequests, "C3_Complaint_Type"},
		{"service request created date", capetown.ServiceServiceRequests, capetown.LayerServiceRequests, "Created_On_Date"},
		{"building plan ward", capetown.ServiceBuildingPlans, capetown.LayerBuildingPlans, "Ward_No"},
		{"building plan submission date", capetown.ServiceBuildingPlans, capetown.LayerBuildingPlans, "Submission_Date"},
		{"building plan approval date", capetown.ServiceBuildingPlans, capetown.LayerBuildingPlans, "Approval_Date"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := liveClient(tc.service)
			info, err := retry(t, tc.name, func(c2 context.Context) (*arcgis.LayerInfo, error) {
				return c.LayerInfo(c2, tc.layerID)
			})
			if err != nil {
				t.Fatalf("LayerInfo(%s/%d): %v", tc.service, tc.layerID, err)
			}
			for _, f := range info.Fields {
				if f.Name == tc.field {
					return
				}
			}
			t.Errorf("field %q not found on %s layer %d (%s)", tc.field, tc.service, tc.layerID, info.Name)
		})
	}
}

// TestLiveHubItemsResolve asserts each ArcGIS Online item ID still points at a
// live FeatureServer, and flags when the City has republished it under a new
// URL so the fallback in ServiceURL can be refreshed.
func TestLiveHubItemsResolve(t *testing.T) {
	for _, svc := range capetown.HubServices() {
		t.Run(svc, func(t *testing.T) {
			id, _ := capetown.HubItemID(svc)
			u, err := retry(t, svc, func(c context.Context) (string, error) {
				return arcgis.ItemURL(c, arcgis.ArcGISOnline, id, arcgis.WithTimeout(liveTimeout))
			})
			if err != nil {
				t.Fatalf("ItemURL(%s): %v", id, err)
			}
			if u != capetown.ServiceURL(svc) {
				t.Logf("%s now resolves to %s; refresh its fallback URL", svc, u)
			}
			c := arcgis.NewClient(u, arcgis.WithTimeout(liveTimeout))
			if _, err := retry(t, svc, c.ServiceInfo); err != nil {
				t.Errorf("ServiceInfo(%s): %v", u, err)
			}
		})
	}
}

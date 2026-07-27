package capetown_test

import (
	"strings"
	"testing"

	capetown "github.com/richardwooding/capetown-opendata"
)

func TestServiceURL(t *testing.T) {
	u := capetown.ServiceURL("ODP_SPLIT_5")
	if !strings.HasSuffix(u, "/FeatureServer") {
		t.Errorf("ServiceURL should point at a FeatureServer, got %q", u)
	}
	if !strings.Contains(u, "/ODP_SPLIT_5/") {
		t.Errorf("ServiceURL should embed the service name, got %q", u)
	}
}

func TestServicesNonEmpty(t *testing.T) {
	s := capetown.Services()
	if len(s) == 0 {
		t.Fatal("Services() returned no services")
	}
	for _, name := range s {
		if !strings.HasPrefix(name, "ODP_SPLIT_") {
			t.Errorf("unexpected service name %q", name)
		}
	}
}

func TestLoadSheddingBlocks(t *testing.T) {
	q := capetown.LoadSheddingBlocks()
	if q.Service != capetown.ServiceLoadShedding {
		t.Errorf("Service = %q, want %q", q.Service, capetown.ServiceLoadShedding)
	}
	if q.Params.LayerID != capetown.LayerLoadSheddingBlocks {
		t.Errorf("LayerID = %d, want %d", q.Params.LayerID, capetown.LayerLoadSheddingBlocks)
	}
	// Pre-built queries no longer pin a field list; the full schema is returned.
	if len(q.Params.Fields) != 0 {
		t.Errorf("expected no pinned fields, got %v", q.Params.Fields)
	}
}

func TestLandParcelsBySuburb(t *testing.T) {
	q := capetown.LandParcelsBySuburb("Newlands")
	if q.Service != capetown.ServiceLandParcels {
		t.Errorf("Service = %q, want %q", q.Service, capetown.ServiceLandParcels)
	}
	if q.Params.LayerID != capetown.LayerLandParcels {
		t.Errorf("LayerID = %d, want %d", q.Params.LayerID, capetown.LayerLandParcels)
	}
	if !strings.Contains(q.Params.Where, "Newlands") {
		t.Errorf("Where = %q, want it to reference Newlands", q.Params.Where)
	}
}

func TestLandParcelsBySuburbEscapesQuotes(t *testing.T) {
	q := capetown.LandParcelsBySuburb("O'Hara")
	if !strings.Contains(q.Params.Where, "O''Hara") {
		t.Errorf("Where = %q, want escaped single quote", q.Params.Where)
	}
}

func TestWaterQualityResultsOrdered(t *testing.T) {
	q := capetown.WaterQualityResults()
	if q.Service != capetown.ServiceWaterQuality {
		t.Errorf("Service = %q, want %q", q.Service, capetown.ServiceWaterQuality)
	}
	if q.Params.LayerID != capetown.LayerWaterQuality {
		t.Errorf("LayerID = %d, want %d", q.Params.LayerID, capetown.LayerWaterQuality)
	}
	if len(q.Params.OrderByFields) == 0 {
		t.Error("expected water quality results to be ordered by date")
	}
}

func TestNamedQueriesHaveLocators(t *testing.T) {
	cases := map[string]capetown.Query{
		"LoadSheddingBlocks": capetown.LoadSheddingBlocks(),
		"Wards":              capetown.Wards(),
		"LandParcels":        capetown.LandParcels(),
		"TaxiRoutes":         capetown.TaxiRoutes(),
		"PublicLighting":     capetown.PublicLighting(),
		"HeritageInventory":  capetown.HeritageInventory(),
		"WaterQuality":       capetown.WaterQualityResults(),
	}
	for name, q := range cases {
		if q.Service == "" {
			t.Errorf("%s: empty service", name)
		}
		if q.Params.LayerID < 0 {
			t.Errorf("%s: layer ID = %d, want >= 0", name, q.Params.LayerID)
		}
	}
}

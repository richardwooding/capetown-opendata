package capetown_test

import (
	"strings"
	"testing"

	capetown "github.com/richardwooding/capetown-opendata"
)

func TestBaseURL(t *testing.T) {
	if !strings.HasSuffix(capetown.BaseURL, "/FeatureServer") {
		t.Errorf("BaseURL should point at a FeatureServer, got %q", capetown.BaseURL)
	}
}

func TestLoadSheddingBlocks(t *testing.T) {
	p := capetown.LoadSheddingBlocks()
	if p.LayerID != capetown.LayerLoadSheddingBlocks {
		t.Errorf("LayerID = %d, want %d", p.LayerID, capetown.LayerLoadSheddingBlocks)
	}
	// Pre-built queries no longer pin a field list; the full schema is returned.
	if len(p.Fields) != 0 {
		t.Errorf("expected no pinned fields, got %v", p.Fields)
	}
}

func TestLandParcelsBySuburb(t *testing.T) {
	p := capetown.LandParcelsBySuburb("Newlands")
	if p.LayerID != capetown.LayerLandParcels {
		t.Errorf("LayerID = %d, want %d", p.LayerID, capetown.LayerLandParcels)
	}
	if !strings.Contains(p.Where, "Newlands") {
		t.Errorf("Where = %q, want it to reference Newlands", p.Where)
	}
}

func TestLandParcelsBySuburbEscapesQuotes(t *testing.T) {
	p := capetown.LandParcelsBySuburb("O'Hara")
	if !strings.Contains(p.Where, "O''Hara") {
		t.Errorf("Where = %q, want escaped single quote", p.Where)
	}
}

func TestWaterQualityResultsOrdered(t *testing.T) {
	p := capetown.WaterQualityResults()
	if p.LayerID != capetown.LayerWaterQuality {
		t.Errorf("LayerID = %d, want %d", p.LayerID, capetown.LayerWaterQuality)
	}
	if len(p.OrderByFields) == 0 {
		t.Error("expected water quality results to be ordered by date")
	}
}

func TestNamedQueriesHaveLayerIDs(t *testing.T) {
	cases := map[string]int{
		"LoadSheddingBlocks":  capetown.LoadSheddingBlocks().LayerID,
		"Wards":               capetown.Wards().LayerID,
		"LandParcels":         capetown.LandParcels().LayerID,
		"TaxiRoutes":          capetown.TaxiRoutes().LayerID,
		"WaterQualityResults": capetown.WaterQualityResults().LayerID,
	}
	for name, id := range cases {
		if id <= 0 {
			t.Errorf("%s: layer ID = %d, want positive", name, id)
		}
	}
}

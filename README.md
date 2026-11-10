# capetown-opendata

[![Go Reference](https://pkg.go.dev/badge/github.com/richardwooding/capetown-opendata.svg)](https://pkg.go.dev/github.com/richardwooding/capetown-opendata)
[![Go](https://github.com/richardwooding/capetown-opendata/actions/workflows/go.yml/badge.svg)](https://github.com/richardwooding/capetown-opendata/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Named dataset locators and pre-built queries for the **City of Cape Town Open
Data Portal**, on top of [`go-arcgis`](https://github.com/richardwooding/go-arcgis).

```go
import (
    arcgis "github.com/richardwooding/go-arcgis"
    capetown "github.com/richardwooding/capetown-opendata"
)

// A dataset is a split service plus a layer ID within it.
q := capetown.LandParcelsBySuburb("Newlands")
client := arcgis.NewClient(capetown.ServiceURL(q.Service))

features, err := client.Layer(q.Params.LayerID).Query().
    From(q.Params).
    WithinEnvelope(18.4, -34.0, 18.6, -33.8).
    All(ctx)
```

## Install

```sh
go get github.com/richardwooding/capetown-opendata
```

Requires Go 1.27+. Depends only on
[`go-arcgis`](https://github.com/richardwooding/go-arcgis).

## What's here

The City publishes its open data across thirteen themed feature services
(`ODP_SPLIT_1` … `ODP_SPLIT_13`), each numbering its layers from zero, plus a
handful of tables hosted on ArcGIS Online. This package names the ones worth
knowing:

- `Services()` lists the split services; `ServiceURL(service)` and
  `ServiceURLIn(folder, service)` build their FeatureServer endpoints on the
  default `esapqa` server or the alternate `citymaps` one.
- `Service*` and `Layer*` constants locate each dataset, and constructors
  return a `Query{Service, Params}` with sensible ordering: `Wards`,
  `OfficialSuburbs` / `OfficialSuburbByName`, `LandParcels` /
  `LandParcelsBySuburb`, `LoadSheddingBlocks`, `TaxiRoutes`, `PublicLighting`,
  `HeritageInventory`, `WaterQualityResults`.
- ArcGIS Online tables are pinned to their hub item IDs, because the City
  republishes them under date-stamped service names: `ServiceRequests`,
  `BuildingPlanApprovals`, `PopulationEstimates` /
  `PopulationEstimateForSuburb` and `TaxiFares` / `TaxiFaresFromRank`.
  `HubServices()` lists their keys and `HubItemID(service)` gives the item ID;
  resolve the current URL with `arcgis.ItemURL` (or `arcgis.ItemInfo` for the
  title and modified date) and fall back to `ServiceURL`.

Constructors do not pin a field list, so each query returns the layer's full
schema; select columns at the call site. Name filters match case-insensitively
because several City layers store names in upper case.

## A note on layer IDs

Service names and layer IDs are best-effort and may drift as the portal is
republished. A daily live check runs the integration tests
(`go test -tags=integration ./...`) against the City's services to catch
drift. Confirm a suspect locator yourself with `Client.ServiceInfo`:

```go
info, _ := arcgis.NewClient(capetown.ServiceURL("ODP_SPLIT_5")).ServiceInfo(ctx)
for _, l := range info.Layers {
    fmt.Println(l.ID, l.Name)
}
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## Sponsor

If this saves you time, you can [sponsor its maintenance](https://github.com/sponsors/richardwooding).
Sponsorship pays for the unglamorous half — triage, dependency bumps, release plumbing — and is
never a condition of getting help here.

## License

[MIT](LICENSE)

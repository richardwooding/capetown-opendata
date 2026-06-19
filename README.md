# capetown-opendata

[![Go Reference](https://pkg.go.dev/badge/github.com/richardwooding/capetown-opendata.svg)](https://pkg.go.dev/github.com/richardwooding/capetown-opendata)
[![Go](https://github.com/richardwooding/capetown-opendata/actions/workflows/go.yml/badge.svg)](https://github.com/richardwooding/capetown-opendata/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Named layer IDs and pre-built queries for the **City of Cape Town Open Data
Portal**, on top of [`go-arcgis`](https://github.com/richardwooding/go-arcgis).

```go
import (
    arcgis "github.com/richardwooding/go-arcgis"
    capetown "github.com/richardwooding/capetown-opendata"
)

client := arcgis.NewClient(capetown.BaseURL)

// Pre-built, optionally refined
features, err := client.Layer(capetown.LayerLoadSheddingBlocks).Query().
    From(capetown.LoadSheddingBlocksForStage(4)).
    WithinEnvelope(18.4, -34.0, 18.6, -33.8).
    All(ctx)
```

## Install

```sh
go get github.com/richardwooding/capetown-opendata
```

Requires Go 1.26+. Depends only on
[`go-arcgis`](https://github.com/richardwooding/go-arcgis).

## What's here

- `BaseURL` — the CCT Open Data Feature Service endpoint.
- Named layer constants (`LayerLoadSheddingBlocks`, `LayerWards`,
  `LayerLandParcels`, `LayerTaxiRoutes`, `LayerWaterQuality`, …).
- Pre-built `arcgis.QueryParams` constructors with sensible default fields and
  ordering: `LoadSheddingBlocks`, `LoadSheddingBlocksForStage`,
  `ServiceRequests`, `ServiceRequestsBySuburb`, `Wards`, `LandParcels`,
  `TaxiRoutes`, `WaterQualityResults`.

Each constructor returns a plain `arcgis.QueryParams`, so you can run it
directly or refine it through the `go-arcgis` fluent builder via `.From(...)`.

## A note on layer IDs

The layer IDs here are best-effort and may drift as the upstream service is
republished. Confirm them against the live service with `Client.ServiceInfo`
if accuracy matters:

```go
info, _ := arcgis.NewClient(capetown.BaseURL).ServiceInfo(ctx)
for _, l := range info.Layers {
    fmt.Println(l.ID, l.Name)
}
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)

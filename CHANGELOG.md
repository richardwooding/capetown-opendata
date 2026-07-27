# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-07-27

The City of Cape Town retired the monolithic `Theme_Based/Open_Data_Service`
and split its layers across a dozen themed feature services (`ODP_SPLIT_1` …
`ODP_SPLIT_12`), renumbering layer IDs per service. Every dataset moved, so the
endpoint model is reworked and callers must migrate.

### Changed (breaking)
- Replaced the single `BaseURL` constant with `BaseFolder` and a
  `ServiceURL(service)` helper — a dataset is now addressed by BOTH a split
  service name and a layer ID.
- Dataset constructors (`Wards`, `LandParcels`, `LandParcelsBySuburb`,
  `LoadSheddingBlocks`, `TaxiRoutes`, `PublicLighting`, `HeritageInventory`,
  `WaterQualityResults`) now return a `Query{Service, Params}` locator instead
  of a bare `arcgis.QueryParams`.
- Remapped every `Layer*` constant to its new per-service ID and added matching
  `Service*` constants for the hosting split service.
- `LandParcelsBySuburb` now matches case-insensitively (`UPPER(...)`), because
  the upstream suburb column stores names in upper case.

### Added
- `Query` locator type, `Service*` constants, and `Services()` listing the
  canonical split services for callers that enumerate the whole portal.

## [0.1.0] - 2026-06-19

Initial release. Extracted from the `capetown` subpackage of
[`go-arcgis`](https://github.com/richardwooding/go-arcgis) into a standalone
module.

### Added
- `BaseURL` for the City of Cape Town Open Data Feature Service.
- Named layer constants for well-known CCT datasets.
- Pre-built `arcgis.QueryParams` constructors: `LoadSheddingBlocks`,
  `LoadSheddingBlocksForStage`, `ServiceRequests`, `ServiceRequestsBySuburb`,
  `Wards`, `LandParcels`, `TaxiRoutes`, and `WaterQualityResults`.

[0.2.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.2.0
[0.1.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.1.0

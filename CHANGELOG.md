# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-10-05

### Added
- Locators for two datasets the City hosts on ArcGIS Online rather than the
  ODP_SPLIT services: `ServiceRequests()` (citizen service requests since 2023)
  and `BuildingPlanApprovals()` (building plan applications since 2014), with
  `ServiceServiceRequests`, `ServiceBuildingPlans` and matching `Layer*`
  constants. Both are non-spatial tables.
- `HubItemID(service)` and `HubServices()`. The hub datasets' service names
  embed dates and go stale, so each is pinned to its item ID; resolve the
  current URL with `arcgis.ItemURL`. `ServiceURL` returns a known fallback URL
  for these keys. `Services()` still lists only the ODP_SPLIT services.

### Changed
- Requires `go-arcgis` v0.4.0.

## [0.2.1] - 2026-10-05

### Changed
- Bumped `go-arcgis` to v0.3.0 (no API change upstream).
- Raised the minimum Go version to 1.27 and applied the Go 1.27 modernizers.

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

[0.3.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.3.0
[0.2.1]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.2.1
[0.2.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.2.0
[0.1.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.1.0

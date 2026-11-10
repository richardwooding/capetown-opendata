# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2026-11-10

### Added
- `ServiceRequestsSpatial()`, with `ServiceServiceRequestsSpatial` and
  `LayerServiceRequestsSpatial`: the City's geocoded subset of the service
  requests table, published as a point layer under hub item
  `6e8f77f27ea84763abbc11023561a406` (about 3.3 million of 4.5 million
  notifications). It takes spatial filters, which the table cannot. Note its
  date columns keep the older names `Created_On_Date`, `Changed_on` and
  `Completed_Date`. `HubServices()` now lists five keys.

## [0.5.0] - 2026-11-10

### Added
- Locators for two suburb- and transport-level datasets the City loaded on
  6 November 2026, both non-spatial tables on ArcGIS Online pinned to their hub
  item IDs: `PopulationEstimates()` / `PopulationEstimateForSuburb(name)`
  (modelled 2025 population and household estimates for 190 suburbs, key
  `ServicePopulationEstimates`) and `TaxiFares()` / `TaxiFaresFromRank(ref)`
  (minibus taxi fares per origin-destination rank pair surveyed 2022-2025, key
  `ServiceTaxiFares`), with matching `Layer*` constants. `HubServices()` now
  lists four keys.
- `OfficialSuburbs()` and `OfficialSuburbByName(name)` on `ODP_SPLIT_5` layer 3,
  with `ServiceOfficialSuburbs` and `LayerOfficialSuburbs`. The population
  table is keyed by suburb name, and this is the polygon it joins to
  (case-insensitively; `OFC_SBRB_NAME` is upper case).

### Changed
- `ServiceRequests()` now points at the City's replacement hub item
  `da3b6d46845e4671bde500e8de7e1d3e` ("Service Requests Jan 2023 to Present").
  The City froze the previous item on 6 October 2026 and marked it deprecated.
  The new table has a reworked schema: `CreatedOnDT`, `ChangedOnDT` and
  `CompletedDT` replace `Created_On_Date` and friends, and `Suburb_Clean`,
  `Longitude` and `Latitude` are new. `Ward` is still a zero-padded string and
  the default `ObjectId DESC` ordering still yields the newest requests first.
- Requires `go-arcgis` v0.6.0.

## [0.4.0] - 2026-10-07

### Changed
- **The default server is now `esapqa.capetown.gov.za`** instead of
  `citymaps.capetown.gov.za`. `BaseFolder` equals the new `FolderESAPQA`.
  The City's open data hub links its datasets to esapqa, which also publishes
  ODP_SPLIT_9 as a feature service and adds ODP_SPLIT_13. Record counts differ
  between the two servers on some layers. Use `FolderCityMaps` with
  `ServiceURLIn` to keep reading from citymaps.

### Added
- `FolderESAPQA` and `FolderCityMaps` constants, and
  `ServiceURLIn(folder, service)` to address a service on any server folder.
- `ODP_SPLIT_13` in `Services()`: beaches, tidal pools, public toilets, bath
  houses, spray parks, stadia, recreation centres, ratepayers' associations,
  dam levels and rainfall. It exists on esapqa only.
- Live tests honour a `CAPETOWN_BASE_FOLDER` environment variable to check a
  server other than the default.

## [0.3.1] - 2026-10-05

### Fixed
- `ServiceRequests()` now orders by `ObjectId DESC` instead of
  `Created_On_Date DESC`. Sorting the five-million-row table by date took
  ArcGIS Online about 20 seconds per page and timed out when paging; the
  object ID is indexed, and its highest values are the newest requests.

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

[0.6.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.6.0
[0.5.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.5.0
[0.4.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.4.0
[0.3.1]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.3.1
[0.3.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.3.0
[0.2.1]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.2.1
[0.2.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.2.0
[0.1.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.1.0

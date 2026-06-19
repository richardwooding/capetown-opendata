# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.0]: https://github.com/richardwooding/capetown-opendata/releases/tag/v0.1.0

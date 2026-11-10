# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A thin Go library providing named layer IDs and pre-built `arcgis.QueryParams`
constructors for the City of Cape Town Open Data Portal, layered on top of
[`go-arcgis`](https://github.com/richardwooding/go-arcgis). It performs no HTTP
itself — callers feed the returned `QueryParams` to a `go-arcgis` client. All
production code lives in the single root file `capetown.go`.

## Commands

```sh
go build ./...
go test -race ./...                          # full suite (matches CI)
go test -race -run TestLoadSheddingBlocks    # single test
golangci-lint run                            # v2 config in .golangci.yml
```

Requires Go 1.27+. CI (`.github/workflows/go.yml`) runs build, race tests, and
golangci-lint v2 on push/PR to `main`.

## Architecture & conventions

- **Three kinds of exports per dataset:** a `Service*` string constant (the
  `ODP_SPLIT_*` service, or a key such as `SERVICE_REQUESTS` for an ArcGIS
  Online dataset), a `Layer*` integer constant (the layer ID within that
  service) and one or more constructors returning a `Query{Service, Params}`
  with a sensible default `OrderByFields`. Constructors never pin `Fields`
  (a test enforces this): the upstream schemas drift, so callers select
  columns themselves.
- **Filtered variants delegate to the base constructor**, then set `.Where`
  (e.g. `LandParcelsBySuburb` calls `LandParcels()` and adds a clause via
  `equalsFold`). Follow this pattern when adding filters rather than
  rebuilding the params from scratch. Match names case-insensitively: several
  City layers store them in upper case.
- **ArcGIS Online datasets** (`hubDatasets`) are pinned to a hub item ID with a
  fallback URL, because the City republishes them under date-stamped service
  names. Add the key to `HubServices()` as well; the tests loop over it.
- **Layer IDs are best-effort** and may drift as the upstream service is
  republished. The unit tests do not touch the network; the daily
  `go test -tags=integration` run (`live-check.yml`) catches drift. Confirm a
  suspect ID against the live service via `arcgis.Client.ServiceInfo` before
  changing a constant.
- **Shared field names** used across datasets (e.g.
  `fieldOfficialSuburbName = "OFC_SBRB_NAME"`) are unexported constants —
  `goconst` requires literals repeated 3+ times to be hoisted.

## Linting notes (golangci-lint v2)

- `revive`'s `exported` rule is on: every exported identifier needs a doc
  comment starting with its name.
- `misspell` uses US locale but **`councillor` is whitelisted** — it is the
  standard South African spelling in CCT data; don't "correct" it.
- `goimports` enforces a local-prefix import group for
  `github.com/richardwooding/capetown-opendata`; formatting via `gofumpt`.
- Test files are exempted from `gocyclo`, `errcheck`, `dupl`, `gosec`,
  `gocritic`, and `goconst`.

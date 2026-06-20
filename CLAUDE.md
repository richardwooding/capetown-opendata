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

Requires Go 1.26+. CI (`.github/workflows/go.yml`) runs build, race tests, and
golangci-lint v2 on push/PR to `main`.

## Architecture & conventions

- **Two kinds of exports per dataset:** a `Layer*` integer constant (the
  ArcGIS layer ID) plus one or more constructor functions returning
  `arcgis.QueryParams` with sensible default `Fields` / `OrderByFields`.
- **Filtered variants delegate to the base constructor**, then set `.Where`
  (e.g. `LoadSheddingBlocksForStage` calls `LoadSheddingBlocks()` and adds a
  `STAGE = N` clause). Follow this pattern when adding filters rather than
  rebuilding the params from scratch.
- **Layer IDs are best-effort** and may drift as the upstream service is
  republished (`LayerServiceRequests` is an explicit placeholder). They are not
  validated against the live service in tests. Confirm a suspect ID against the
  live service via `arcgis.Client.ServiceInfo` before changing a constant.
- **Shared field names** used across datasets (e.g. `fieldSuburb = "SUBURB"`)
  are unexported constants — `goconst` requires literals repeated 3+ times to
  be hoisted.

## Linting notes (golangci-lint v2)

- `revive`'s `exported` rule is on: every exported identifier needs a doc
  comment starting with its name.
- `misspell` uses US locale but **`councillor` is whitelisted** — it is the
  standard South African spelling in CCT data; don't "correct" it.
- `goimports` enforces a local-prefix import group for
  `github.com/richardwooding/capetown-opendata`; formatting via `gofumpt`.
- Test files are exempted from `gocyclo`, `errcheck`, `dupl`, `gosec`,
  `gocritic`, and `goconst`.

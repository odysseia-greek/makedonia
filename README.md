# Makedonia

A cohesive toolkit of small, focused services for Greek lexicography and text analysis. This mono‑repo powers dictionary lookups, phrase and fuzzy search, analytics, and GraphQL aggregation.

At the edge sits Alexandros — the GraphQL gateway — which orchestrates specialized backend services over gRPC. Shared protocol definitions live in Filippos so every service speaks the same language.

## Why this repo exists
- Consistent developer experience across all services (shared Make targets, docs generation, proto lifecycle)
- Clear separation of concerns: each service does one thing well (exact, partial, phrase, fuzzy, extended retrieval, analytics, seeding, etc.)
- Single source of truth for APIs and documentation

## High‑level architecture
- Alexandros (GraphQL) aggregates and fans out requests to domain services via gRPC
- Domain services implement different retrieval strategies (exact, substring, phrase, fuzzy, extended texts)
- Demokritos seeds and scaffolds dictionaries and texts
- Eukleides captures usage analytics/metrics
- Dareios exercises the system end‑to‑end for confidence
- Filippos provides shared protobuf contracts used by all the above

```
Client ➜ Alexandros (GraphQL) ➜ gRPC calls ➜ Antigonos | Hefaistion | Perdikkas | Parmenion | Ptolemaios | Eukleides
                                                       ↳ Demokritos (data seeding)
                                                       ↳ Dareios (system tests)
```

## Services (at a glance)

- Alexandros — GraphQL gateway coordinating dictionary/search services
- Demokritos — Seeder for initial data setup and dictionary scaffolding
- Hefaistion — Exact‑match search
- Perdikkas — Partial (substring) search
- Antigonos — Fuzzy (approximate) search
- Parmenion — Phrase/extended search
- Ptolemaios — Extended results retrieval; library for texts
- Dareios — System testing (end‑to‑end/integration)
- Eukleides — Analytics and user metrics (e.g., top searches, usage patterns)
- Filippos — Shared proto layer for gRPC services

## Repository layout
- alexandros/ … GraphQL gateway, schema, middleware, routing, SpectaQL docs
- antigonos/, hefaistion/, perdikkas/, parmenion/, ptolemaios/ … search services + proto and generated code
- demokritos/ … data seeding and corpus scaffolding
- eukleides/ … analytics/metrics collection and proto
- dareios/ … e2e and integration tests helpers
- filippos/ … shared protobuf definitions and generated artifacts
- Makefile … common developer tasks (codegen, docs, images)

## Quick start

Prerequisites
- Go toolchain (matching the versions declared per module)
- buf (for protobuf code generation): https://buf.build
- Docker or Podman (optional) for container images
- spectaql (optional) to regenerate the GraphQL HTML docs

Common tasks
- Tidy and format all modules:
  ```bash
  make tidy
  ```
- Generate protobuf code for all gRPC services:
  ```bash
  make generate
  ```
- Generate API docs (gRPC + GraphQL SpectaQL):
  ```bash
  make docs
  ```

Build container images (dev/prod tags)
- Bump and build images using the helper script:
  ```bash
  # DEV images
  OWNER="your-registry-or-user" ROOT="github.com/odysseia-greek/makedonia" make images-dev

  # PROD images
  OWNER="your-registry-or-user" ROOT="github.com/odysseia-greek/makedonia" make images-prod
  ```
  The above delegates to `./bump-images.sh` and uses `Containerfile` definitions in each service.

Running services
- Each service provides its own `Containerfile` and module. You can:
  - build and run containers with Docker/Podman, or
  - run from source via `go run` within the specific module.
- See the `docs/` folder inside each service for details, ports, and examples where available.

## Documentation
- GraphQL reference (Alexandros): generated with SpectaQL in `alexandros/docs`
- gRPC API docs: generated per service into their respective `docs/` directories
- Protobuf definitions: centralized under `filippos/` and referenced by domain services

## Versioning
This project follows semantic versioning per module. Convenience scripts exist at the repo root:
- `bump-patch.sh` — bump patch versions across modules
- `bump-minor.sh` — bump minor versions across modules

## Contributing
- Open an issue or discussion describing the change
- For proto changes: update `filippos` first, run `make generate`, and regenerate docs with `make docs`
- Keep changes focused per service and add/update docs when behavior changes

## License
Licensed under the terms of the LICENSE file in this repository.

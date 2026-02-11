# Dareios – Alexandros GraphQL Ginkgo Suite

This directory contains a Ginkgo v2 test suite that exercises the Alexandros GraphQL server.

## Configure endpoint

Set the environment variable `ALEXANDROS_URL` to the GraphQL HTTP endpoint of Alexandros. If not set, the suite defaults to `http://localhost:8080/query`.

Example:

```
export ALEXANDROS_URL=http://localhost:8080/query
```

## Run the suite

Using `go test`:

```
go test ./dareios/alexandros -v
```

Or using Ginkgo CLI (if installed):

```
ginkgo -r dareios/alexandros -v
```

## Covered queries

- `health` – asserts server and service health structure
- `counterTopFive` – asserts list shape and non-negative counts
- `fuzzy` – runs a small search and validates paging/result structure

The tests use lightweight local response structs to avoid tight coupling with the Alexandros code generation. If you prefer, you can adapt them to import types from `alexandros/graph/model` later.

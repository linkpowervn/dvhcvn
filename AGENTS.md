# DVHCVN - Vietnamese Administrative Units

Go library (`go 1.24`, zero external dependencies) for fetching and querying Vietnam's province/district/ward data from a remote JSON API.

## Package structure

```
dvhcvn/         # main library package (import: github.com/linkpowervn/dvhcvn/v1)
  dvhcvn.go     # types + Service + API methods
  dvhcvn_test.go
example/
  main.go       # runnable demo using the library
```

Notable: the `dvhcvn` package lives in a subdirectory that matches the package name. Import paths use the full prefix `github.com/linkpowervn/dvhcvn/v1`.

## Key types

- `Level1` → province (fields: `Level1ID`, `Name`, `Type`, `Level2s`)
- `Level2` → district (`Level2ID`, `Name`, `Type`, `Level3s`)
- `Level3` → ward (`Level3ID`, `Name`, `Type`)
- `DVHCResponse` → wraps API response as `{"data": [...]}`
- `Service` → holds `remotePath` string; no HTTP client abstraction

## API methods (all on `*Service`)

| Method | Returns |
|---|---|
| `GetProvinces()` | `[]Level1, error` |
| `GetProvince(id)` | `*Level1, error` |
| `GetDistricts(provinceID)` | `[]Level2, error` |
| `GetDistrict(provinceID, districtID)` | `*Level2, error` |
| `GetWards(provinceID, districtID)` | `[]Level3, error` |
| `GetWard(provinceID, districtID, wardID)` | `*Level3, error` |

Non-existent IDs return `fmt.Errorf` (not custom error type). Each lookup method walks the nested data in memory rather than making separate API calls.

## Test quirks

- **Tests require network access.** They call a real remote endpoint (`https://public-assets.hcm.s3storage.vn/json/dvhcvn/v1.json`). They will fail offline.
- No mocks, no test fixtures, no `httptest.Server`. The HTTP call is hardcoded inside `Service.fetchData()` using `http.Get` — not injected or mockable.
- Tests do not clean up state between runs (there is none — `Service` is stateless).
- Running `go test -v -count=1 -timeout 30s ./dvhcvn` runs all 8 tests (~2.5s with good connectivity).
- `TestFetchData` is an integration test that duplicates `TestGetProvinces` coverage.
- Benchmark tests exist (`BenchmarkGetProvinces`, `BenchmarkFetchData`) — run with `go test -bench=. ./dvhcvn`.

## Commands

```sh
go build ./...              # builds library + example
go vet ./...                # static analysis (no violations)
go test -v -count=1 ./dvhcvn  # run tests (requires network)
go test -bench=. ./dvhcvn     # run benchmarks (requires network)
go run ./example              # run the demo program (requires network)
```

No Makefile, no CI config, no linters beyond `go vet`.

## Remote data source

The JSON endpoint is an S3-hosted file (~63 provinces, each with nested districts/wards). Data structure:

```json
{"data": [{ "level1_id": "01", "name": "...", "type": "...", "level2s": [...] }]}
```

The source-of-truth upstream is `github.com/daohoangson/dvhcvn`.

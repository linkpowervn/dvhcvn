# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About

`github.com/linkpowervn/dvhcvn` là thư viện Go thuần (không có external dependencies) để truy vấn dữ liệu đơn vị hành chính Việt Nam 3 cấp (tỉnh/thành → huyện/quận → xã/phường) từ một remote JSON endpoint.

> **v2 đang phát triển** tại thư mục `v2/` (`github.com/linkpowervn/dvhcvn/v2`). v2 hỗ trợ mô hình 2 cấp (tỉnh/thành → phường/xã/đặc khu) theo cải cách hành chính 7/2025, dùng dataset từ `thanglequoc/vietnamese-provinces-database`. v1 ở repo root được giữ nguyên để không phá vỡ người dùng hiện tại.

## Commands

```bash
# Build
go build ./...

# Chạy toàn bộ tests (yêu cầu kết nối internet — tests gọi API thực)
go test ./...

# Chạy một test cụ thể
go test ./v1/... -run TestGetProvinces

# Chạy benchmarks
go test ./v1/... -bench=.

# Chạy example
go run example/main.go
```

## Architecture

Toàn bộ logic nằm trong package `v1/dvhcvn.go`:

- **`Service`** — wrapper quanh một remote URL. Mỗi lần gọi method đều fetch lại toàn bộ JSON (không cache).
- **Hierarchy**: `Level1` (tỉnh/thành phố) → `Level2` (quận/huyện) → `Level3` (phường/xã). Dữ liệu được lồng nhau trong JSON.
- **`fetchData()`** — parse `DVHCResponse{Data: []Level1}` từ response. Các method cấp cao hơn (`GetDistricts`, `GetWards`, v.v.) đều gọi qua chain này, nghĩa là mỗi query tốn một HTTP round-trip.
- Data source thực tế: `https://public-assets.hcm.s3storage.vn/json/dvhcvn/dvhcvn.json`

## Notes

- Tests v1 là integration tests gọi HTTP thực — không mock network.
- Không có caching ở v1: `GetWard(p, d, w)` thực hiện 1 HTTP fetch → linear scan qua toàn bộ data.

## v2 (`v2/`)

- Module path: `github.com/linkpowervn/dvhcvn/v2`
- Go commands phải chạy từ trong thư mục `v2/` (module riêng biệt)
- Tests offline hoàn toàn — dùng `httptest.NewServer` + fixtures trong `v2/testdata/`
- Ba variant JSON được hỗ trợ: `vn_only_simplified`, `simplified`, `full` — auto-detect từ nội dung JSON
- Cache lazy: fetch 1 lần, gọi `Refresh(ctx)` để làm mới

```bash
# Build v2
cd v2 && go build ./...

# Test v2 offline (mặc định, không cần internet)
cd v2 && go test -race ./...

# Test integration thủ công (cần internet)
cd v2 && go test -run TestIntegration ./...

# Benchmarks v2
cd v2 && go test -bench=. -benchmem ./...
```

# DVHCVN - Đơn vị hành chính Việt Nam

Thư viện Go để truy vấn dữ liệu đơn vị hành chính Việt Nam.

---

## v2 (khuyến nghị) — mô hình 2 cấp, cải cách 7/2025

**[→ Xem README v2](v2/README.md)**

```bash
go get github.com/linkpowervn/dvhcvn/v2
```

v2 hỗ trợ cấu trúc hành chính mới: **tỉnh/thành phố → phường/xã/đặc khu** (không còn quận/huyện), dữ liệu từ [thanglequoc/vietnamese-provinces-database](https://github.com/thanglequoc/vietnamese-provinces-database). Có cache lazy, `context.Context`, functional options và validation schema.

---

## v1 (legacy) — mô hình 3 cấp

```bash
go get github.com/linkpowervn/dvhcvn
```

v1 hỗ trợ cấu trúc cũ: **tỉnh/thành phố → quận/huyện → phường/xã**, dữ liệu từ repo [daohoangson/dvhcvn](https://github.com/daohoangson/dvhcvn). Không có cache — mỗi lần gọi đều fetch lại JSON từ remote.

```go
import "github.com/linkpowervn/dvhcvn/v1"

svc := dvhcvn.NewService("https://public-assets.hcm.s3storage.vn/json/dvhcvn/dvhcvn.json")

provinces, _ := svc.GetProvinces()
districts, _ := svc.GetDistricts("56")
wards, _     := svc.GetWards("56", "568")
ward, _      := svc.GetWard("56", "568", "22363")
```

**v1 không còn được phát triển thêm.** Khuyến nghị migrate sang v2 — xem [hướng dẫn migrate](v2/README.md#migrate-từ-v1).

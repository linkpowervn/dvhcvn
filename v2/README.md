# dvhcvn/v2 — Đơn vị hành chính Việt Nam

Thư viện Go thuần (không external dependency) để truy vấn dữ liệu đơn vị hành chính Việt Nam **2 cấp** (tỉnh/thành phố → phường/xã/đặc khu) theo cải cách hành chính tháng 7/2025.

Nguồn dữ liệu: [thanglequoc/vietnamese-provinces-database](https://github.com/thanglequoc/vietnamese-provinces-database) — MIT license, cập nhật theo các Nghị quyết của Quốc hội và UBTVQH.

> **Đang dùng v1?** Xem [hướng dẫn migrate](#migrate-từ-v1).

## Cài đặt

```bash
go get github.com/linkpowervn/dvhcvn/v2
```

## Sử dụng cơ bản

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/linkpowervn/dvhcvn/v2"
)

func main() {
    svc := dvhcvn.NewService(dvhcvn.DefaultDatasetURL)
    ctx := context.Background()

    // Lấy tất cả tỉnh/thành phố
    provinces, err := svc.Provinces(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tổng số tỉnh/thành: %d\n", len(provinces))

    // Lấy một tỉnh theo mã
    hanoi, err := svc.Province(ctx, "01")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(hanoi.FullName) // Thành phố Hà Nội

    // Lấy tất cả phường/xã của một tỉnh
    wards, err := svc.WardsByProvinceCode(ctx, "01")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Số phường/xã Hà Nội: %d\n", len(wards))

    // Lấy một phường/xã theo mã
    ward, err := svc.Ward(ctx, "00004")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ward.FullName) // Phường Ba Đình
}
```

## Cấu hình

```go
import "net/http"

svc := dvhcvn.NewService(
    dvhcvn.DefaultDatasetURL,
    dvhcvn.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
    dvhcvn.WithUserAgent("my-app/1.0"),
)
```

Dùng URL tuỳ chỉnh (ví dụ mirror nội bộ hoặc variant JSON khác):

```go
const myURL = "https://cdn.jsdelivr.net/gh/thanglequoc/vietnamese-provinces-database@v3.1.0/json/full_json_generated_data_vn_units.json"
svc := dvhcvn.NewService(myURL)
```

## Cache và Refresh

Dữ liệu được fetch **một lần duy nhất** khi có query đầu tiên, sau đó cache trong memory suốt vòng đời của `Service`. Gọi `Refresh` để tải lại dữ liệu mới nhất:

```go
if err := svc.Refresh(ctx); err != nil {
    // fetch thất bại — cache cũ vẫn được giữ nguyên
    log.Println("refresh failed:", err)
}
```

## Tìm kiếm

```go
// Tìm tỉnh/thành theo tên (không phân biệt hoa/thường, hỗ trợ tiếng Anh)
results, _ := svc.FindProvincesByName(ctx, "minh")
// → [Thành phố Hồ Chí Minh]

// Tìm phường/xã theo tên
wards, _ := svc.FindWardsByName(ctx, "ba dinh")
// → [Phường Ba Đình, ...]

// Tra cứu phường/xã theo slug (CodeName)
wards, _ := svc.FindWardsByCodeName(ctx, "ba_dinh")
```

> **Lưu ý:** Tìm kiếm so sánh sau khi `strings.ToLower`, không chuẩn hoá dấu tiếng Việt. "ba dinh" khớp với "Ba Dinh" (NameEn) nhưng không khớp với "Ba Đình" (Name). Để tìm theo tiếng Việt có dấu, truyền nguyên dấu: `"Ba Đình"`.

## Biến thể JSON (Variants)

Dataset upstream có 3 biến thể, thư viện **tự nhận diện** từ nội dung JSON:

| Variant | Fields có thêm | Dùng khi |
|---|---|---|
| `vn_only_simplified` | `Code`, `FullName` | Chỉ cần tên tiếng Việt, tối giản |
| `simplified` | + `Name`, `NameEn`, `FullName`, `FullNameEn`, `CodeName` | Đa ngôn ngữ, slug cho URL |
| `full` | + `Type`, `AdministrativeUnit*` | Cần phân biệt phường/xã/đặc khu |

`DefaultDatasetURL` trỏ tới biến thể `simplified` (pin tag `v3.1.0`).

Kiểm tra variant sau khi load:

```go
svc.Provinces(ctx) // trigger load
fmt.Println(svc.DetectedVariant()) // "simplified"
```

Chỉ định variant thủ công (thêm validation chặt hơn):

```go
svc := dvhcvn.NewService(url, dvhcvn.WithVariant(dvhcvn.VariantFull))
```

## Cấu trúc dữ liệu

```go
type Province struct {
    Code       string // "01"
    Name       string // "Hà Nội"           (simplified, full)
    NameEn     string // "Ha Noi"            (simplified, full)
    FullName   string // "Thành phố Hà Nội"
    FullNameEn string // "Ha Noi City"       (simplified, full)
    CodeName   string // "ha_noi"            (simplified, full)
    Type       string // "province"          (full)
    // AdministrativeUnit* fields            (full)
    Wards []Ward
}

type Ward struct {
    Code         string // "00004"
    Name         string // "Ba Đình"
    NameEn       string // "Ba Dinh"
    FullName     string // "Phường Ba Đình"
    FullNameEn   string // "Ba Dinh Ward"
    CodeName     string // "ba_dinh"
    ProvinceCode string // "01"
    Type         string // "ward" / "commune" / "special_administrative_region" (full)
    // AdministrativeUnit* fields (full)
}
```

Mã code là chuỗi zero-padded: tỉnh 2 chữ số (`"01"`–`"99"`), phường/xã 5 chữ số (`"00004"`).

## API

```go
// Khởi tạo
dvhcvn.NewService(url string, opts ...Option) *Service

// Cấu hình
dvhcvn.WithHTTPClient(c *http.Client) Option
dvhcvn.WithVariant(v Variant) Option
dvhcvn.WithUserAgent(ua string) Option

// Lifecycle
svc.Refresh(ctx context.Context) error
svc.DetectedVariant() Variant

// Tra cứu
svc.Provinces(ctx) ([]Province, error)
svc.Province(ctx, code) (Province, error)
svc.Wards(ctx) ([]Ward, error)
svc.Ward(ctx, code) (Ward, error)
svc.WardsByProvinceCode(ctx, provinceCode) ([]Ward, error)

// Tìm kiếm
svc.FindProvincesByName(ctx, substr) ([]Province, error)
svc.FindWardsByName(ctx, substr) ([]Ward, error)
svc.FindWardsByCodeName(ctx, codeName) ([]Ward, error)
```

## Xử lý lỗi

```go
var ve *dvhcvn.ValidationError

switch {
case errors.Is(err, dvhcvn.ErrNotFound):
    // Không tìm thấy theo code/slug
case errors.Is(err, dvhcvn.ErrHTTP):
    // HTTP request thất bại (timeout, status != 200, v.v.)
case errors.As(err, &ve):
    // Dữ liệu JSON không hợp lệ — ve.Path, ve.Field, ve.Reason có chi tiết
}
```

## Migrate từ v1

| v1 | v2 |
|---|---|
| `import "github.com/linkpowervn/dvhcvn/v1"` | `import "github.com/linkpowervn/dvhcvn/v2"` |
| `GetProvinces()` | `Provinces(ctx)` |
| `GetProvince(id)` | `Province(ctx, code)` |
| `GetDistricts(provID)` | ~~Không còn~~ — cấp quận/huyện đã bị bỏ |
| `GetDistrict(provID, distID)` | ~~Không còn~~ |
| `GetWards(provID, distID)` | `WardsByProvinceCode(ctx, provCode)` |
| `GetWard(provID, distID, wardID)` | `Ward(ctx, wardCode)` — ward code là globally unique |
| `Level1ID` / `Level2ID` / `Level3ID` | `Code` |
| Fetch mỗi lần gọi | Cache lazy, dùng `Refresh(ctx)` để làm mới |
| Không có `context.Context` | `context.Context` trên mọi method |

**Thay đổi data model quan trọng:** Không còn cấp quận/huyện (`Level2`). Phường/xã giờ treo thẳng dưới tỉnh/thành phố và có field `ProvinceCode` để biết tỉnh cha. Nếu code cũ dùng quận/huyện để lọc, hãy chuyển sang dùng `WardsByProvinceCode`.

## Phát triển

```bash
# Build
cd v2 && go build ./...

# Test offline (không cần internet)
cd v2 && go test -race ./...

# Test integration (cần internet, hit upstream thật)
cd v2 && go test -run TestIntegration ./...

# Benchmarks
cd v2 && go test -bench=. -benchmem ./...
```

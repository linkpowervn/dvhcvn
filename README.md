# DVHCVN - Đơn vị hành chính Việt Nam

Thư viện Go để xử lý dữ liệu đơn vị hành chính Việt Nam (tỉnh, huyện, xã) từ remote API.

Repo chứa data tại đây: [dvhcvn](https://github.com/daohoangson/dvhcvn)

## Cài đặt

```bash
go get github.com/linkpowervn/dvhcvn
```

## Sử dụng

```go
package main

import (
    "fmt"
    "log"
    "github.com/linkpowervn/dvhcvn/dvhcvn"
)

func main() {
    // Khởi tạo service với URL API
    service := dvhcvn.NewService("https://api.example.com/dvhc")
    
    // Lấy danh sách tỉnh
    provinces, err := service.GetProvinces()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, province := range provinces {
        fmt.Printf("Tỉnh: %s (ID: %s)\n", province.Name, province.Level1ID)
    }
    
    // Lấy thông tin một tỉnh cụ thể
    province, err := service.GetProvince("56")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tỉnh: %s\n", province.Name)
    
    // Lấy danh sách huyện của tỉnh
    districts, err := service.GetDistricts("56")
    if err != nil {
        log.Fatal(err)
    }
    
    for _, district := range districts {
        fmt.Printf("Huyện: %s (ID: %s)\n", district.Name, district.Level2ID)
    }
    
    // Lấy thông tin một huyện cụ thể
    district, err := service.GetDistrict("56", "568")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Huyện: %s\n", district.Name)
    
    // Lấy danh sách xã/phường của huyện
    wards, err := service.GetWards("56", "568")
    if err != nil {
        log.Fatal(err)
    }
    
    for _, ward := range wards {
        fmt.Printf("Xã/Phường: %s (ID: %s)\n", ward.Name, ward.Level3ID)
    }
    
    // Lấy thông tin một xã/phường cụ thể
    ward, err := service.GetWard("56", "568", "22363")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Xã/Phường: %s\n", ward.Name)
}
```

## Cấu trúc dữ liệu

Thư viện hỗ trợ cấu trúc dữ liệu JSON như sau:

```json
{
  "level1_id": "56",
  "name": "Tỉnh Khánh Hòa",
  "type": "Tỉnh",
  "level2s": [
    {
      "level2_id": "568",
      "name": "Thành phố Nha Trang",
      "type": "Thành phố",
      "level3s": [
        {
          "level3_id": "22363",
          "name": "Phường Lộc Thọ",
          "type": "Phường"
        }
      ]
    }
  ]
}
```

## API Methods

- `GetProvinces()` - Lấy danh sách tất cả tỉnh/thành phố
- `GetProvince(provinceID)` - Lấy thông tin một tỉnh/thành phố cụ thể
- `GetDistricts(provinceID)` - Lấy danh sách huyện/quận của một tỉnh
- `GetDistrict(provinceID, districtID)` - Lấy thông tin một huyện/quận cụ thể
- `GetWards(provinceID, districtID)` - Lấy danh sách xã/phường của một huyện
- `GetWard(provinceID, districtID, wardID)` - Lấy thông tin một xã/phường cụ thể


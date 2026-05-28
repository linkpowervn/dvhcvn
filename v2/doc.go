// Package dvhcvn cung cấp API truy vấn đơn vị hành chính Việt Nam 2 cấp
// (tỉnh/thành phố → phường/xã/đặc khu) theo cải cách hành chính tháng 7/2025.
//
// Dữ liệu lấy từ https://github.com/thanglequoc/vietnamese-provinces-database
// và hỗ trợ cả 3 biến thể JSON: vn_only_simplified, simplified, full.
//
// Sử dụng cơ bản:
//
//	svc := dvhcvn.NewService(dvhcvn.DefaultDatasetURL)
//	provinces, err := svc.Provinces(ctx)
//	wards, err := svc.WardsByProvinceCode(ctx, "01") // Hà Nội
package dvhcvn

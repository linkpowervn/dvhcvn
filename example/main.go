package main

import (
	"fmt"
	"log"

	"github.com/linkpowervn/dvhcvn/dvhcvn"
)

func main() {
	// Khởi tạo service với URL API thực tế
	service := dvhcvn.NewService("https://public-assets.hcm.s3storage.vn/json/dvhcvn/dvhcvn.json")

	// Lấy danh sách tỉnh
	fmt.Println("=== Danh sách tỉnh ===")
	provinces, err := service.GetProvinces()
	if err != nil {
		log.Printf("Lỗi khi lấy danh sách tỉnh: %v", err)
		return
	}

	for _, province := range provinces {
		fmt.Printf("Tỉnh: %s (ID: %s, Loại: %s)\n", province.Name, province.Level1ID, province.Type)
	}

	// Lấy thông tin một tỉnh cụ thể (Hà Nội)
	fmt.Println("\n=== Thông tin Thành phố Hà Nội ===")
	province, err := service.GetProvince("01")
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin tỉnh: %v", err)
		return
	}
	fmt.Printf("Tỉnh: %s (ID: %s, Loại: %s)\n", province.Name, province.Level1ID, province.Type)

	// Lấy danh sách huyện của tỉnh
	fmt.Println("\n=== Danh sách huyện/quận ===")
	districts, err := service.GetDistricts("01")
	if err != nil {
		log.Printf("Lỗi khi lấy danh sách huyện: %v", err)
		return
	}

	for i, district := range districts {
		fmt.Printf("%d. %s (ID: %s, Loại: %s)\n", i+1, district.Name, district.Level2ID, district.Type)
		if i >= 4 { // Chỉ hiển thị 5 huyện đầu tiên
			fmt.Printf("... và %d huyện/quận khác\n", len(districts)-5)
			break
		}
	}

	// Lấy thông tin một huyện cụ thể (Quận Ba Đình)
	fmt.Println("\n=== Thông tin Quận Ba Đình ===")
	district, err := service.GetDistrict("01", "001")
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin huyện: %v", err)
		return
	}
	fmt.Printf("Huyện: %s (ID: %s, Loại: %s)\n", district.Name, district.Level2ID, district.Type)

	// Lấy danh sách xã/phường của huyện
	fmt.Println("\n=== Danh sách xã/phường ===")
	wards, err := service.GetWards("01", "001")
	if err != nil {
		log.Printf("Lỗi khi lấy danh sách xã/phường: %v", err)
		return
	}

	for i, ward := range wards {
		fmt.Printf("%d. %s (ID: %s, Loại: %s)\n", i+1, ward.Name, ward.Level3ID, ward.Type)
		if i >= 4 { // Chỉ hiển thị 5 xã/phường đầu tiên
			fmt.Printf("... và %d xã/phường khác\n", len(wards)-5)
			break
		}
	}

	// Lấy thông tin một xã/phường cụ thể (Phường Phúc Xá)
	fmt.Println("\n=== Thông tin Phường Phúc Xá ===")
	ward, err := service.GetWard("01", "001", "00001")
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin xã/phường: %v", err)
		return
	}
	fmt.Printf("Xã/Phường: %s (ID: %s, Loại: %s)\n", ward.Name, ward.Level3ID, ward.Type)

	fmt.Println("\n=== Hoàn thành demo ===")
}
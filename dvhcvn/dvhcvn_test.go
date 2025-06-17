package dvhcvn

import (
	"testing"
)

const testRemotePath = "https://public-assets.hcm.s3storage.vn/json/dvhcvn/dvhcvn.json"

func TestNewService(t *testing.T) {
	service := NewService(testRemotePath)
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	if service.remotePath != testRemotePath {
		t.Errorf("Expected remotePath %s, got %s", testRemotePath, service.remotePath)
	}
}

func TestGetProvinces(t *testing.T) {
	service := NewService(testRemotePath)
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Error("Expected at least one province")
	}
	
	// Kiểm tra province đầu tiên có dữ liệu hợp lệ
	province := provinces[0]
	if province.Level1ID == "" {
		t.Error("Province Level1ID is empty")
	}
	if province.Name == "" {
		t.Error("Province Name is empty")
	}
	if province.Type == "" {
		t.Error("Province Type is empty")
	}
	
	t.Logf("Found province: %s (ID: %s, Type: %s)", province.Name, province.Level1ID, province.Type)
}

func TestGetProvince(t *testing.T) {
	service := NewService(testRemotePath)
	
	// Lấy danh sách provinces trước để có ID hợp lệ
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Fatal("No provinces found")
	}
	
	provinceID := provinces[0].Level1ID
	province, err := service.GetProvince(provinceID)
	if err != nil {
		t.Fatalf("GetProvince failed: %v", err)
	}
	if province.Level1ID != provinceID {
		t.Errorf("Expected province ID %s, got %s", provinceID, province.Level1ID)
	}
	
	t.Logf("Found province: %s (ID: %s)", province.Name, province.Level1ID)
	
	// Test với ID không tồn tại
	_, err = service.GetProvince("999999")
	if err == nil {
		t.Error("Expected error for non-existent province ID")
	}
}

func TestGetDistricts(t *testing.T) {
	service := NewService(testRemotePath)
	
	// Lấy province ID hợp lệ
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Fatal("No provinces found")
	}
	
	provinceID := provinces[0].Level1ID
	districts, err := service.GetDistricts(provinceID)
	if err != nil {
		t.Fatalf("GetDistricts failed: %v", err)
	}
	
	if len(districts) == 0 {
		t.Error("Expected at least one district")
	}
	
	// Kiểm tra district đầu tiên
	district := districts[0]
	if district.Level2ID == "" {
		t.Error("District Level2ID is empty")
	}
	if district.Name == "" {
		t.Error("District Name is empty")
	}
	
	t.Logf("Found %d districts in province %s", len(districts), provinceID)
	t.Logf("First district: %s (ID: %s)", district.Name, district.Level2ID)
}

func TestGetDistrict(t *testing.T) {
	service := NewService(testRemotePath)
	
	// Lấy province và district ID hợp lệ
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Fatal("No provinces found")
	}
	
	provinceID := provinces[0].Level1ID
	districts, err := service.GetDistricts(provinceID)
	if err != nil {
		t.Fatalf("GetDistricts failed: %v", err)
	}
	if len(districts) == 0 {
		t.Fatal("No districts found")
	}
	
	districtID := districts[0].Level2ID
	district, err := service.GetDistrict(provinceID, districtID)
	if err != nil {
		t.Fatalf("GetDistrict failed: %v", err)
	}
	if district.Level2ID != districtID {
		t.Errorf("Expected district ID %s, got %s", districtID, district.Level2ID)
	}
	
	t.Logf("Found district: %s (ID: %s)", district.Name, district.Level2ID)
	
	// Test với ID không tồn tại
	_, err = service.GetDistrict(provinceID, "999999")
	if err == nil {
		t.Error("Expected error for non-existent district ID")
	}
}

func TestGetWards(t *testing.T) {
	service := NewService(testRemotePath)
	
	// Lấy province và district ID hợp lệ
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Fatal("No provinces found")
	}
	
	provinceID := provinces[0].Level1ID
	districts, err := service.GetDistricts(provinceID)
	if err != nil {
		t.Fatalf("GetDistricts failed: %v", err)
	}
	if len(districts) == 0 {
		t.Fatal("No districts found")
	}
	
	districtID := districts[0].Level2ID
	wards, err := service.GetWards(provinceID, districtID)
	if err != nil {
		t.Fatalf("GetWards failed: %v", err)
	}
	
	if len(wards) == 0 {
		t.Error("Expected at least one ward")
	}
	
	// Kiểm tra ward đầu tiên
	ward := wards[0]
	if ward.Level3ID == "" {
		t.Error("Ward Level3ID is empty")
	}
	if ward.Name == "" {
		t.Error("Ward Name is empty")
	}
	
	t.Logf("Found %d wards in district %s", len(wards), districtID)
	t.Logf("First ward: %s (ID: %s)", ward.Name, ward.Level3ID)
}

func TestGetWard(t *testing.T) {
	service := NewService(testRemotePath)
	
	// Lấy province, district và ward ID hợp lệ
	provinces, err := service.GetProvinces()
	if err != nil {
		t.Fatalf("GetProvinces failed: %v", err)
	}
	if len(provinces) == 0 {
		t.Fatal("No provinces found")
	}
	
	provinceID := provinces[0].Level1ID
	districts, err := service.GetDistricts(provinceID)
	if err != nil {
		t.Fatalf("GetDistricts failed: %v", err)
	}
	if len(districts) == 0 {
		t.Fatal("No districts found")
	}
	
	districtID := districts[0].Level2ID
	wards, err := service.GetWards(provinceID, districtID)
	if err != nil {
		t.Fatalf("GetWards failed: %v", err)
	}
	if len(wards) == 0 {
		t.Fatal("No wards found")
	}
	
	wardID := wards[0].Level3ID
	ward, err := service.GetWard(provinceID, districtID, wardID)
	if err != nil {
		t.Fatalf("GetWard failed: %v", err)
	}
	if ward.Level3ID != wardID {
		t.Errorf("Expected ward ID %s, got %s", wardID, ward.Level3ID)
	}
	
	t.Logf("Found ward: %s (ID: %s)", ward.Name, ward.Level3ID)
	
	// Test với ID không tồn tại
	_, err = service.GetWard(provinceID, districtID, "999999")
	if err == nil {
		t.Error("Expected error for non-existent ward ID")
	}
}

func TestFetchData(t *testing.T) {
	service := NewService(testRemotePath)
	data, err := service.fetchData()
	if err != nil {
		t.Fatalf("fetchData failed: %v", err)
	}
	if data == nil {
		t.Fatal("fetchData returned nil data")
	}
	if len(data) == 0 {
		t.Fatal("fetchData returned empty data")
	}
	if data[0].Level1ID == "" {
		t.Error("First province Level1ID is empty")
	}
	if data[0].Name == "" {
		t.Error("First province Name is empty")
	}
	
	t.Logf("Fetched %d provinces, first: %s (ID: %s)", len(data), data[0].Name, data[0].Level1ID)
}

// Benchmark tests
func BenchmarkGetProvinces(b *testing.B) {
	service := NewService(testRemotePath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetProvinces()
		if err != nil {
			b.Fatalf("GetProvinces failed: %v", err)
		}
	}
}

func BenchmarkFetchData(b *testing.B) {
	service := NewService(testRemotePath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.fetchData()
		if err != nil {
			b.Fatalf("fetchData failed: %v", err)
		}
	}
}
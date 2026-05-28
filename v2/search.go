package dvhcvn

import (
	"context"
	"strings"
)

// FindProvincesByName tìm kiếm tỉnh/thành phố theo substring (không phân biệt hoa/thường).
// Tìm trong các field: Name, NameEn, FullName, FullNameEn.
func (s *Service) FindProvincesByName(ctx context.Context, substr string) ([]Province, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(substr)
	var result []Province
	for _, p := range snap.provinces {
		if containsLower(p.Name, lower) || containsLower(p.NameEn, lower) ||
			containsLower(p.FullName, lower) || containsLower(p.FullNameEn, lower) {
			result = append(result, p)
		}
	}
	return result, nil
}

// FindWardsByName tìm kiếm phường/xã theo substring (không phân biệt hoa/thường).
// Tìm trong các field: Name, NameEn, FullName, FullNameEn.
func (s *Service) FindWardsByName(ctx context.Context, substr string) ([]Ward, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(substr)
	var result []Ward
	for _, w := range snap.wards {
		if containsLower(w.Name, lower) || containsLower(w.NameEn, lower) ||
			containsLower(w.FullName, lower) || containsLower(w.FullNameEn, lower) {
			result = append(result, w)
		}
	}
	return result, nil
}

// FindWardsByCodeName tìm tất cả phường/xã có CodeName khớp chính xác với slug.
// CodeName là slug dạng snake_case (ví dụ: "ba_dinh", "ngoc_ha").
// Lưu ý: CodeName không đảm bảo unique toàn cục; có thể có nhiều kết quả.
func (s *Service) FindWardsByCodeName(ctx context.Context, codeName string) ([]Ward, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	var result []Ward
	for _, w := range snap.wards {
		if w.CodeName == codeName {
			result = append(result, w)
		}
	}
	return result, nil
}

func containsLower(s, lowerSubstr string) bool {
	return strings.Contains(strings.ToLower(s), lowerSubstr)
}

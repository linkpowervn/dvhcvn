package dvhcvn

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// DefaultDatasetURL là URL mặc định trỏ tới dataset thanglequoc (simplified, CDN, pin tag v3.1.0).
const DefaultDatasetURL = "https://cdn.jsdelivr.net/gh/thanglequoc/vietnamese-provinces-database@v3.1.0/json/simplified_json_generated_data_vn_units_minified.json"

// Service là entry point để truy vấn dữ liệu đơn vị hành chính.
// Lần đầu gọi bất kỳ method nào sẽ fetch và parse JSON, sau đó cache trong memory.
// Gọi Refresh để tải lại dữ liệu.
type Service struct {
	cfg    config
	url    string
	snap   atomic.Pointer[snapshot]
	loadMu sync.Mutex
}

// NewService tạo Service mới với URL và các tuỳ chọn cấu hình.
func NewService(url string, opts ...Option) *Service {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return &Service{cfg: cfg, url: url}
}

// DetectedVariant trả về variant đã nhận diện sau lần load đầu tiên.
// Trả về chuỗi rỗng nếu chưa load.
func (s *Service) DetectedVariant() Variant {
	snap := s.snap.Load()
	if snap == nil {
		return ""
	}
	return snap.variant
}

// Refresh tải lại toàn bộ dữ liệu từ remote và cập nhật cache.
// Nếu request thất bại, cache cũ vẫn được giữ nguyên.
func (s *Service) Refresh(ctx context.Context) error {
	s.loadMu.Lock()
	defer s.loadMu.Unlock()
	return s.load(ctx)
}

// ensure đảm bảo snapshot đã được load (lazy). An toàn cho nhiều goroutine.
func (s *Service) ensure(ctx context.Context) (*snapshot, error) {
	if snap := s.snap.Load(); snap != nil {
		return snap, nil
	}
	s.loadMu.Lock()
	defer s.loadMu.Unlock()
	if snap := s.snap.Load(); snap != nil {
		return snap, nil
	}
	if err := s.load(ctx); err != nil {
		return nil, err
	}
	return s.snap.Load(), nil
}

// load fetch + parse + validate + build snapshot. Phải gọi với loadMu held.
func (s *Service) load(ctx context.Context) error {
	provinces, err := s.fetchAndParse(ctx)
	if err != nil {
		return err
	}

	variant := s.cfg.variant
	if variant == VariantAuto {
		variant = detectVariant(provinces)
	}

	if err := validate(provinces, variant); err != nil {
		return err
	}

	s.snap.Store(buildSnapshot(provinces, variant))
	return nil
}

// Provinces trả về danh sách tất cả tỉnh/thành phố.
func (s *Service) Provinces(ctx context.Context) ([]Province, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	return snap.provinces, nil
}

// Province trả về tỉnh/thành phố theo mã code.
// Trả về ErrNotFound nếu không tồn tại.
func (s *Service) Province(ctx context.Context, code string) (Province, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return Province{}, err
	}
	p, ok := snap.provinceByCode[code]
	if !ok {
		return Province{}, fmt.Errorf("province %q: %w", code, ErrNotFound)
	}
	return *p, nil
}

// Wards trả về danh sách phẳng tất cả phường/xã/đặc khu.
func (s *Service) Wards(ctx context.Context) ([]Ward, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	return snap.wards, nil
}

// Ward trả về phường/xã/đặc khu theo mã code.
// Trả về ErrNotFound nếu không tồn tại.
func (s *Service) Ward(ctx context.Context, code string) (Ward, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return Ward{}, err
	}
	w, ok := snap.wardByCode[code]
	if !ok {
		return Ward{}, fmt.Errorf("ward %q: %w", code, ErrNotFound)
	}
	return *w, nil
}

// WardsByProvinceCode trả về tất cả phường/xã của một tỉnh/thành phố.
// Trả về ErrNotFound nếu mã tỉnh không tồn tại.
func (s *Service) WardsByProvinceCode(ctx context.Context, provinceCode string) ([]Ward, error) {
	snap, err := s.ensure(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := snap.provinceByCode[provinceCode]; !ok {
		return nil, fmt.Errorf("province %q: %w", provinceCode, ErrNotFound)
	}
	ws := snap.wardsByProv[provinceCode]
	result := make([]Ward, len(ws))
	for i, w := range ws {
		result[i] = *w
	}
	return result, nil
}

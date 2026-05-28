package dvhcvn_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linkpowervn/dvhcvn/v2"
)

func serveFixture(t *testing.T, path string) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	var count atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	t.Cleanup(srv.Close)
	return srv, &count
}

func serveStatus(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// ── Variant detection ─────────────────────────────────────────────────────────

func TestDetectVnOnly(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/vn_only.json")
	svc := dvhcvn.NewService(srv.URL)

	provinces, err := svc.Provinces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(provinces) != 2 {
		t.Fatalf("want 2 provinces, got %d", len(provinces))
	}
	if svc.DetectedVariant() != dvhcvn.VariantVnOnlySimplified {
		t.Errorf("want VariantVnOnlySimplified, got %q", svc.DetectedVariant())
	}
	if provinces[0].Name != "" {
		t.Error("vn_only should not have Name field")
	}
}

func TestDetectSimplified(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	if _, err := svc.Provinces(context.Background()); err != nil {
		t.Fatal(err)
	}
	if svc.DetectedVariant() != dvhcvn.VariantSimplified {
		t.Errorf("want VariantSimplified, got %q", svc.DetectedVariant())
	}
}

func TestDetectFull(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/full.json")
	svc := dvhcvn.NewService(srv.URL)

	provinces, err := svc.Provinces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if svc.DetectedVariant() != dvhcvn.VariantFull {
		t.Errorf("want VariantFull, got %q", svc.DetectedVariant())
	}
	if provinces[0].AdministrativeUnitID == nil {
		t.Error("full variant should have AdministrativeUnitID")
	}
}

func TestDetectedVariantBeforeLoad(t *testing.T) {
	svc := dvhcvn.NewService("http://localhost")
	if svc.DetectedVariant() != "" {
		t.Errorf("want empty before load, got %q", svc.DetectedVariant())
	}
}

// ── Core lookups ──────────────────────────────────────────────────────────────

func TestProvinces(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	provinces, err := svc.Provinces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(provinces) != 2 {
		t.Fatalf("want 2, got %d", len(provinces))
	}
	if provinces[0].Code != "01" {
		t.Errorf("want Code 01, got %q", provinces[0].Code)
	}
	if provinces[0].Name != "Hà Nội" {
		t.Errorf("unexpected Name: %q", provinces[0].Name)
	}
}

func TestProvince(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()

	p, err := svc.Province(ctx, "79")
	if err != nil {
		t.Fatal(err)
	}
	if p.FullName != "Thành phố Hồ Chí Minh" {
		t.Errorf("unexpected FullName: %q", p.FullName)
	}
}

func TestProvinceNotFound(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.Province(context.Background(), "99")
	if !errors.Is(err, dvhcvn.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestWards(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	wards, err := svc.Wards(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(wards) != 4 {
		t.Fatalf("want 4 wards total, got %d", len(wards))
	}
}

func TestWard(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	w, err := svc.Ward(context.Background(), "00004")
	if err != nil {
		t.Fatal(err)
	}
	if w.Name != "Ba Đình" {
		t.Errorf("unexpected Name: %q", w.Name)
	}
	if w.ProvinceCode != "01" {
		t.Errorf("unexpected ProvinceCode: %q", w.ProvinceCode)
	}
}

func TestWardNotFound(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.Ward(context.Background(), "99999")
	if !errors.Is(err, dvhcvn.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestWardsByProvinceCode(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	wards, err := svc.WardsByProvinceCode(context.Background(), "01")
	if err != nil {
		t.Fatal(err)
	}
	if len(wards) != 2 {
		t.Fatalf("want 2 wards for Hà Nội, got %d", len(wards))
	}
}

func TestWardsByProvinceCodeNotFound(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.WardsByProvinceCode(context.Background(), "99")
	if !errors.Is(err, dvhcvn.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

// ── Validation ────────────────────────────────────────────────────────────────

func TestValidationDuplicateCode(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/invalid_dup_codes.json")
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.Provinces(context.Background())
	if !errors.Is(err, dvhcvn.ErrInvalidSchema) {
		t.Errorf("want ErrInvalidSchema, got %v", err)
	}
	var ve *dvhcvn.ValidationError
	if !errors.As(err, &ve) {
		t.Errorf("want *ValidationError, got %T", err)
	}
	if ve.Field != "Code" {
		t.Errorf("want Field=Code, got %q", ve.Field)
	}
}

func TestValidationOrphanWard(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/invalid_orphan_ward.json")
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.Provinces(context.Background())
	if !errors.Is(err, dvhcvn.ErrInvalidSchema) {
		t.Errorf("want ErrInvalidSchema, got %v", err)
	}
	var ve *dvhcvn.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	if ve.Field != "ProvinceCode" {
		t.Errorf("want Field=ProvinceCode, got %q", ve.Field)
	}
}

func TestValidationHTTPError(t *testing.T) {
	srv := serveStatus(t, http.StatusInternalServerError)
	svc := dvhcvn.NewService(srv.URL)

	_, err := svc.Provinces(context.Background())
	if !errors.Is(err, dvhcvn.ErrHTTP) {
		t.Errorf("want ErrHTTP, got %v", err)
	}
}

// ── Concurrency (cold cache, single HTTP request) ─────────────────────────────

func TestConcurrentLoad(t *testing.T) {
	srv, count := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			if _, err := svc.Provinces(context.Background()); err != nil {
				t.Errorf("concurrent Provinces error: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := count.Load(); got != 1 {
		t.Errorf("want exactly 1 HTTP request, got %d", got)
	}
}

// ── Refresh ───────────────────────────────────────────────────────────────────

func TestRefreshUpdatesData(t *testing.T) {
	fixture1, err := os.ReadFile("testdata/simplified.json")
	if err != nil {
		t.Fatal(err)
	}
	fixture2, err := os.ReadFile("testdata/vn_only.json")
	if err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	current := fixture1
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		data := current
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	t.Cleanup(srv.Close)

	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()

	if _, err := svc.Provinces(ctx); err != nil {
		t.Fatal(err)
	}
	if svc.DetectedVariant() != dvhcvn.VariantSimplified {
		t.Fatalf("initial variant: want Simplified, got %q", svc.DetectedVariant())
	}

	mu.Lock()
	current = fixture2
	mu.Unlock()

	if err := svc.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
	if svc.DetectedVariant() != dvhcvn.VariantVnOnlySimplified {
		t.Errorf("after refresh: want VnOnlySimplified, got %q", svc.DetectedVariant())
	}
}

func TestRefreshFailureKeepsOldSnapshot(t *testing.T) {
	data, err := os.ReadFile("testdata/simplified.json")
	if err != nil {
		t.Fatal(err)
	}

	var fail atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	t.Cleanup(srv.Close)

	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()

	if _, err := svc.Provinces(ctx); err != nil {
		t.Fatal(err)
	}

	fail.Store(true)
	if err := svc.Refresh(ctx); !errors.Is(err, dvhcvn.ErrHTTP) {
		t.Errorf("want ErrHTTP from Refresh, got %v", err)
	}

	// Snapshot cũ vẫn phải hoạt động
	provinces, err := svc.Provinces(ctx)
	if err != nil {
		t.Fatalf("old snapshot should still work: %v", err)
	}
	if len(provinces) == 0 {
		t.Error("expected non-empty provinces from old snapshot")
	}
}

// ── Custom HTTP client ────────────────────────────────────────────────────────

func TestCustomHTTPClientTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client := &http.Client{Timeout: 10 * time.Millisecond}
	svc := dvhcvn.NewService(srv.URL, dvhcvn.WithHTTPClient(client))

	_, err := svc.Provinces(context.Background())
	if !errors.Is(err, dvhcvn.ErrHTTP) {
		t.Errorf("want ErrHTTP (timeout), got %v", err)
	}
}

// ── Context cancellation ──────────────────────────────────────────────────────

func TestContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	svc := dvhcvn.NewService(srv.URL)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel ngay lập tức

	_, err := svc.Provinces(ctx)
	if !errors.Is(err, dvhcvn.ErrHTTP) {
		t.Errorf("want ErrHTTP (cancelled), got %v", err)
	}
}

// ── Search ────────────────────────────────────────────────────────────────────

func TestFindProvincesByName(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()

	results, err := svc.FindProvincesByName(ctx, "hà nội")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "01" {
		t.Errorf("unexpected results: %+v", results)
	}

	// Case-insensitive
	results2, err := svc.FindProvincesByName(ctx, "HÀ NỘI")
	if err != nil {
		t.Fatal(err)
	}
	if len(results2) != 1 {
		t.Errorf("case-insensitive search failed: %+v", results2)
	}

	// NameEn match
	results3, err := svc.FindProvincesByName(ctx, "ha noi")
	if err != nil {
		t.Fatal(err)
	}
	if len(results3) != 1 {
		t.Errorf("NameEn search failed: %+v", results3)
	}
}

func TestFindWardsByName(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()

	results, err := svc.FindWardsByName(ctx, "bến")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("want 2 wards matching 'bến', got %d", len(results))
	}
}

func TestFindWardsByCodeName(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	results, err := svc.FindWardsByCodeName(context.Background(), "ba_dinh")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "00004" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestFindWardsByCodeNameNotFound(t *testing.T) {
	srv, _ := serveFixture(t, "testdata/simplified.json")
	svc := dvhcvn.NewService(srv.URL)

	results, err := svc.FindWardsByCodeName(context.Background(), "khong_ton_tai")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("want empty, got %+v", results)
	}
}

// ── Benchmarks ────────────────────────────────────────────────────────────────

func BenchmarkProvincesWarmCache(b *testing.B) {
	data, _ := os.ReadFile("testdata/simplified.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer srv.Close()

	svc := dvhcvn.NewService(srv.URL)
	ctx := context.Background()
	svc.Provinces(ctx) // warm

	b.ResetTimer()
	for b.Loop() {
		svc.Provinces(ctx)
	}
}

func BenchmarkColdLoad(b *testing.B) {
	data, _ := os.ReadFile("testdata/simplified.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer srv.Close()

	for b.Loop() {
		svc := dvhcvn.NewService(srv.URL)
		svc.Provinces(context.Background())
	}
}

// ── Integration test (chỉ chạy với -tags integration) ────────────────────────

func TestIntegrationDefaultURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	svc := dvhcvn.NewService(dvhcvn.DefaultDatasetURL)
	ctx := context.Background()

	provinces, err := svc.Provinces(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(provinces) < 30 {
		t.Errorf("expected at least 30 provinces (post-2025 reform), got %d", len(provinces))
	}

	p, err := svc.Province(ctx, "01")
	if err != nil {
		t.Fatal(err)
	}
	if p.FullName == "" {
		t.Error("province 01 FullName should not be empty")
	}
	t.Logf("DetectedVariant: %s, Provinces: %d, Province 01: %s", svc.DetectedVariant(), len(provinces), p.FullName)
}

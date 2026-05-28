package dvhcvn

import "net/http"

type config struct {
	client    *http.Client
	variant   Variant
	userAgent string
}

func defaultConfig() config {
	return config{
		client:    http.DefaultClient,
		variant:   VariantAuto,
		userAgent: "dvhcvn-go/v2",
	}
}

// Option là functional option để cấu hình Service.
type Option func(*config)

// WithHTTPClient dùng *http.Client tuỳ chỉnh (timeout, transport, v.v.).
func WithHTTPClient(c *http.Client) Option {
	return func(cfg *config) { cfg.client = c }
}

// WithVariant chỉ định biến thể JSON cụ thể thay vì auto-detect.
func WithVariant(v Variant) Option {
	return func(cfg *config) { cfg.variant = v }
}

// WithUserAgent đặt User-Agent header cho HTTP request.
func WithUserAgent(ua string) Option {
	return func(cfg *config) { cfg.userAgent = ua }
}

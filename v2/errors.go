package dvhcvn

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound được trả về khi không tìm thấy province hoặc ward theo code/slug.
	ErrNotFound = errors.New("dvhcvn: not found")
	// ErrInvalidSchema được trả về khi dữ liệu JSON không hợp lệ.
	ErrInvalidSchema = errors.New("dvhcvn: invalid schema")
	// ErrHTTP được trả về khi HTTP request thất bại.
	ErrHTTP = errors.New("dvhcvn: http error")
)

// ValidationError chứa thông tin chi tiết về lỗi validation schema.
type ValidationError struct {
	Variant Variant
	Path    string // ví dụ: "provinces[3].wards[12].ProvinceCode"
	Field   string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("dvhcvn: validation error at %s (field %q): %s", e.Path, e.Field, e.Reason)
}

func (e *ValidationError) Unwrap() error { return ErrInvalidSchema }

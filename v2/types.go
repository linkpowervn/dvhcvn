package dvhcvn

// Variant xác định biến thể JSON của dataset thanglequoc/vietnamese-provinces-database.
type Variant string

const (
	// VariantAuto tự động nhận diện biến thể từ nội dung JSON.
	VariantAuto Variant = "auto"
	// VariantVnOnlySimplified chỉ có Code và FullName (tiếng Việt).
	VariantVnOnlySimplified Variant = "vn_only_simplified"
	// VariantSimplified có thêm Name, NameEn, FullNameEn, CodeName.
	VariantSimplified Variant = "simplified"
	// VariantFull có thêm Type và thông tin AdministrativeUnit đầy đủ.
	VariantFull Variant = "full"
)

// Province đại diện cho một tỉnh/thành phố trực thuộc trung ương.
// Các field optional chỉ được điền tuỳ theo Variant của dataset.
type Province struct {
	Code     string `json:"Code"`
	Name     string `json:"Name,omitempty"`
	NameEn   string `json:"NameEn,omitempty"`
	FullName string `json:"FullName"`
	FullNameEn string `json:"FullNameEn,omitempty"`
	CodeName string `json:"CodeName,omitempty"`
	Type     string `json:"Type,omitempty"` // "province" — chỉ có trong VariantFull

	AdministrativeUnitID          *int   `json:"AdministrativeUnitId,omitempty"`
	AdministrativeUnitShortName   string `json:"AdministrativeUnitShortName,omitempty"`
	AdministrativeUnitFullName    string `json:"AdministrativeUnitFullName,omitempty"`
	AdministrativeUnitShortNameEn string `json:"AdministrativeUnitShortNameEn,omitempty"`
	AdministrativeUnitFullNameEn  string `json:"AdministrativeUnitFullNameEn,omitempty"`

	Wards []Ward `json:"Wards"`
}

// Ward đại diện cho một phường/xã/đặc khu hành chính.
type Ward struct {
	Code         string `json:"Code"`
	Name         string `json:"Name,omitempty"`
	NameEn       string `json:"NameEn,omitempty"`
	FullName     string `json:"FullName"`
	FullNameEn   string `json:"FullNameEn,omitempty"`
	CodeName     string `json:"CodeName,omitempty"`
	ProvinceCode string `json:"ProvinceCode"`
	Type         string `json:"Type,omitempty"` // "ward"/"commune"/"special_administrative_region"

	AdministrativeUnitID          *int   `json:"AdministrativeUnitId,omitempty"`
	AdministrativeUnitShortName   string `json:"AdministrativeUnitShortName,omitempty"`
	AdministrativeUnitFullName    string `json:"AdministrativeUnitFullName,omitempty"`
	AdministrativeUnitShortNameEn string `json:"AdministrativeUnitShortNameEn,omitempty"`
	AdministrativeUnitFullNameEn  string `json:"AdministrativeUnitFullNameEn,omitempty"`
}

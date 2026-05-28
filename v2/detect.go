package dvhcvn

// detectVariant nhận diện biến thể JSON dựa trên fields có mặt trong province đầu tiên.
// Thứ tự ưu tiên: VariantFull > VariantSimplified > VariantVnOnlySimplified.
func detectVariant(provinces []Province) Variant {
	if len(provinces) == 0 {
		return VariantVnOnlySimplified
	}
	p := provinces[0]
	if p.AdministrativeUnitID != nil || p.AdministrativeUnitShortName != "" {
		return VariantFull
	}
	if p.Name != "" || p.CodeName != "" {
		return VariantSimplified
	}
	return VariantVnOnlySimplified
}

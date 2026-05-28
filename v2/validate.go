package dvhcvn

import (
	"fmt"
	"strconv"
)

func validate(provinces []Province, variant Variant) error {
	if len(provinces) == 0 {
		return &ValidationError{
			Variant: variant,
			Path:    "provinces",
			Field:   "",
			Reason:  "dataset is empty",
		}
	}

	seenProvinceCode := make(map[string]bool, len(provinces))
	seenWardCode := make(map[string]bool)

	for pi, p := range provinces {
		pPath := fmt.Sprintf("provinces[%d]", pi)

		if err := checkRequired(p, variant, pPath); err != nil {
			return err
		}
		if err := checkProvinceCodeFormat(p.Code, pPath); err != nil {
			return err
		}
		if seenProvinceCode[p.Code] {
			return &ValidationError{
				Variant: variant,
				Path:    pPath,
				Field:   "Code",
				Reason:  fmt.Sprintf("duplicate province code %q", p.Code),
			}
		}
		seenProvinceCode[p.Code] = true

		for wi, w := range p.Wards {
			wPath := fmt.Sprintf("%s.wards[%d]", pPath, wi)

			if err := checkWardRequired(w, variant, wPath); err != nil {
				return err
			}
			if err := checkWardCodeFormat(w.Code, wPath); err != nil {
				return err
			}
			if w.ProvinceCode != p.Code {
				return &ValidationError{
					Variant: variant,
					Path:    wPath,
					Field:   "ProvinceCode",
					Reason:  fmt.Sprintf("ward.ProvinceCode %q does not match parent province.Code %q", w.ProvinceCode, p.Code),
				}
			}
			if seenWardCode[w.Code] {
				return &ValidationError{
					Variant: variant,
					Path:    wPath,
					Field:   "Code",
					Reason:  fmt.Sprintf("duplicate ward code %q", w.Code),
				}
			}
			seenWardCode[w.Code] = true
		}
	}
	return nil
}

func checkRequired(p Province, variant Variant, path string) error {
	if p.Code == "" {
		return &ValidationError{Variant: variant, Path: path, Field: "Code", Reason: "required"}
	}
	if p.FullName == "" {
		return &ValidationError{Variant: variant, Path: path, Field: "FullName", Reason: "required"}
	}
	switch variant {
	case VariantSimplified, VariantFull:
		if p.Name == "" {
			return &ValidationError{Variant: variant, Path: path, Field: "Name", Reason: fmt.Sprintf("required for variant %q", variant)}
		}
		if p.CodeName == "" {
			return &ValidationError{Variant: variant, Path: path, Field: "CodeName", Reason: fmt.Sprintf("required for variant %q", variant)}
		}
	}
	if variant == VariantFull && p.AdministrativeUnitID == nil {
		return &ValidationError{Variant: variant, Path: path, Field: "AdministrativeUnitId", Reason: fmt.Sprintf("required for variant %q", variant)}
	}
	return nil
}

func checkWardRequired(w Ward, variant Variant, path string) error {
	if w.Code == "" {
		return &ValidationError{Variant: variant, Path: path, Field: "Code", Reason: "required"}
	}
	if w.FullName == "" {
		return &ValidationError{Variant: variant, Path: path, Field: "FullName", Reason: "required"}
	}
	if w.ProvinceCode == "" {
		return &ValidationError{Variant: variant, Path: path, Field: "ProvinceCode", Reason: "required"}
	}
	switch variant {
	case VariantSimplified, VariantFull:
		if w.Name == "" {
			return &ValidationError{Variant: variant, Path: path, Field: "Name", Reason: fmt.Sprintf("required for variant %q", variant)}
		}
		if w.CodeName == "" {
			return &ValidationError{Variant: variant, Path: path, Field: "CodeName", Reason: fmt.Sprintf("required for variant %q", variant)}
		}
	}
	return nil
}

// checkProvinceCodeFormat kiểm tra code tỉnh phải là đúng 2 chữ số.
func checkProvinceCodeFormat(code, path string) error {
	if len(code) != 2 {
		return &ValidationError{Path: path, Field: "Code", Reason: fmt.Sprintf("province code must be 2 digits, got %q", code)}
	}
	if _, err := strconv.Atoi(code); err != nil {
		return &ValidationError{Path: path, Field: "Code", Reason: fmt.Sprintf("province code must be numeric, got %q", code)}
	}
	return nil
}

// checkWardCodeFormat kiểm tra code phường/xã phải là đúng 5 chữ số.
func checkWardCodeFormat(code, path string) error {
	if len(code) != 5 {
		return &ValidationError{Path: path, Field: "Code", Reason: fmt.Sprintf("ward code must be 5 digits, got %q", code)}
	}
	if _, err := strconv.Atoi(code); err != nil {
		return &ValidationError{Path: path, Field: "Code", Reason: fmt.Sprintf("ward code must be numeric, got %q", code)}
	}
	return nil
}

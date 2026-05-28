package dvhcvn

type snapshot struct {
	provinces      []Province
	provinceByCode map[string]*Province
	wards          []Ward
	wardByCode     map[string]*Ward
	wardsByProv    map[string][]*Ward
	variant        Variant
}

func buildSnapshot(provinces []Province, variant Variant) *snapshot {
	snap := &snapshot{
		provinces:      provinces,
		provinceByCode: make(map[string]*Province, len(provinces)),
		wardByCode:     make(map[string]*Ward),
		wardsByProv:    make(map[string][]*Ward, len(provinces)),
		variant:        variant,
	}

	for i := range provinces {
		p := &snap.provinces[i]
		snap.provinceByCode[p.Code] = p
		for j := range p.Wards {
			w := &snap.provinces[i].Wards[j]
			snap.wardByCode[w.Code] = w
			snap.wardsByProv[p.Code] = append(snap.wardsByProv[p.Code], w)
			snap.wards = append(snap.wards, *w)
		}
	}

	return snap
}

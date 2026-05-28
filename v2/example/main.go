package main

import (
	"context"
	"fmt"
	"log"

	"github.com/linkpowervn/dvhcvn/v2"
)

func main() {
	svc := dvhcvn.NewService(dvhcvn.DefaultDatasetURL)
	ctx := context.Background()

	provinces, err := svc.Provinces(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Detected variant: %s\n", svc.DetectedVariant())
	fmt.Printf("Total provinces: %d\n\n", len(provinces))

	// In 3 tỉnh đầu tiên
	for _, p := range provinces[:3] {
		fmt.Printf("[%s] %s\n", p.Code, p.FullName)
	}

	// Tra cứu phường/xã của Hà Nội (code "01")
	fmt.Println("\n--- Một số phường/xã của Hà Nội ---")
	wards, err := svc.WardsByProvinceCode(ctx, "01")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tổng số phường/xã Hà Nội: %d\n", len(wards))
	for _, w := range wards[:3] {
		fmt.Printf("  [%s] %s\n", w.Code, w.FullName)
	}

	// Tìm kiếm
	fmt.Println("\n--- Tìm province có 'minh' ---")
	found, err := svc.FindProvincesByName(ctx, "minh")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range found {
		fmt.Printf("  %s\n", p.FullName)
	}
}

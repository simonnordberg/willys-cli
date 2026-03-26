package cmd

import (
	"strings"
	"testing"

	"github.com/simonnordberg/willys-cli/internal/willys"
)

func TestFormatProduct(t *testing.T) {
	p := willys.Product{
		Name:             "Mjölk Färsk 3%",
		Code:             "100010649_ST",
		Price:            "21,90 kr",
		Manufacturer:     "Falköpings",
		DisplayVolume:    "1,5l",
		ComparePrice:     "14,60 kr",
		ComparePriceUnit: "l",
	}
	got := FormatProduct(p)
	// Code first, then price, compare price, description
	for _, want := range []string{"100010649_ST", "21,90 kr", "14,60 kr/l", "Mjölk Färsk 3%", "[Falköpings]", "1,5l"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatProduct missing %q in %q", want, got)
		}
	}
	// Code should appear before description
	codeIdx := strings.Index(got, "100010649_ST")
	nameIdx := strings.Index(got, "Mjölk")
	if codeIdx > nameIdx {
		t.Errorf("code should appear before name: %q", got)
	}
}

func TestFormatProduct_Minimal(t *testing.T) {
	p := willys.Product{
		Name:  "Banan",
		Code:  "100254920_KG",
		Price: "19,90 kr",
	}
	got := FormatProduct(p)
	if !strings.Contains(got, "100254920_KG") || !strings.Contains(got, "Banan") {
		t.Errorf("FormatProduct minimal: %q", got)
	}
	if strings.Contains(got, "[]") {
		t.Errorf("should skip empty manufacturer: %q", got)
	}
}

func TestFormatCart_Empty(t *testing.T) {
	got := FormatCart(willys.Cart{})
	if got != "Cart is empty." {
		t.Errorf("FormatCart empty = %q", got)
	}
}

func TestFormatCart(t *testing.T) {
	cart := willys.Cart{
		TotalUnitCount: 3,
		TotalPrice:     "85,80 kr",
		Products: []willys.CartProduct{
			{
				Name:         "Mjölk",
				Code:         "100010649_ST",
				PickQuantity: 2,
				TotalPrice:   "43,80 kr",
			},
			{
				Name:         "Ägg 15p",
				Code:         "101241403_ST",
				Manufacturer: "Garant Eko",
				PickQuantity: 1,
				TotalPrice:   "62,90 kr",
			},
		},
	}
	got := FormatCart(cart)
	// Code first, then qty, then price
	for _, want := range []string{"Cart — 3 items", "100010649_ST", "101241403_ST", "[Garant Eko]", "Total: 85,80 kr"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatCart missing %q in:\n%s", want, got)
		}
	}
	// Verify code appears before description on same line
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "100010649_ST") {
			codeIdx := strings.Index(line, "100010649_ST")
			nameIdx := strings.Index(line, "Mjölk")
			if codeIdx > nameIdx {
				t.Errorf("code should appear before name: %q", line)
			}
		}
	}
}

func TestFormatOrderHistory_Empty(t *testing.T) {
	got := FormatOrderHistory(nil)
	if got != "No orders found." {
		t.Errorf("FormatOrderHistory empty = %q", got)
	}
}

func TestFormatOrderHistory(t *testing.T) {
	orders := []willys.OrderSummary{
		{
			OrderNumber:  "3057837654",
			DeliveryDate: "2026-03-24",
			OrderStatus:  willys.OrderStatus{Code: "delivered"},
			Total:        "3 291,29 kr",
		},
	}
	got := FormatOrderHistory(orders)
	for _, want := range []string{"1 orders", "3057837654", "2026-03-24", "delivered", "3 291,29 kr"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatOrderHistory missing %q in:\n%s", want, got)
		}
	}
	// No # prefix
	if strings.Contains(got, "#3057837654") {
		t.Errorf("order number should not have # prefix: %s", got)
	}
}

func TestFormatOrderDetail(t *testing.T) {
	order := willys.OrderDetail{
		OrderNumber:   "3057837654",
		StatusDisplay: "Levererad",
		TotalPrice:    willys.FormattedPrice{FormattedValue: "3 291,29 kr"},
		Products: map[string][]willys.OrderEntry{
			"Mejeri": {
				{
					Name:          "Mjölk Färsk 3%",
					Code:          "100010649_ST",
					Manufacturer:  "Falköpings",
					DisplayVolume: "1,5l",
					PickQuantity:  2,
					TotalPrice:    "43,80 kr",
				},
			},
		},
	}
	got := FormatOrderDetail(order)
	// Header: Order NUMBER — STATUS — TOTAL
	if !strings.Contains(got, "Order 3057837654 — Levererad — 3 291,29 kr") {
		t.Errorf("header format wrong in:\n%s", got)
	}
	// Category as plain header
	if !strings.Contains(got, "Mejeri\n") {
		t.Errorf("category should be plain header in:\n%s", got)
	}
	// Item line: code first
	for _, want := range []string{"100010649_ST", "Mjölk Färsk 3%", "[Falköpings]", "1,5l", "43,80 kr"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatOrderDetail missing %q in:\n%s", want, got)
		}
	}
}

func TestFormatCategory(t *testing.T) {
	cat := willys.Category{
		Title: "Frukt & Grönt",
		URL:   "/frukt-gront",
		Children: []willys.Category{
			{Title: "Frukt", URL: "/frukt"},
		},
	}
	got := FormatCategory(cat, 0)
	if !strings.Contains(got, "Frukt & Grönt") || !strings.Contains(got, "/frukt-gront") {
		t.Errorf("root category: %q", got)
	}
	if !strings.Contains(got, "  Frukt") {
		t.Errorf("child should be indented: %q", got)
	}
}

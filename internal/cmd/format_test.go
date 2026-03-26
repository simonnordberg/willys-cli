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
	for _, want := range []string{"Mjölk Färsk 3%", "[Falköpings]", "1,5l", "21,90 kr", "(14,60 kr/l)", "(100010649_ST)"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatProduct missing %q in %q", want, got)
		}
	}
}

func TestFormatProduct_Minimal(t *testing.T) {
	p := willys.Product{
		Name:  "Banan",
		Code:  "100254920_KG",
		Price: "19,90 kr",
	}
	got := FormatProduct(p)
	if !strings.Contains(got, "Banan") || !strings.Contains(got, "(100254920_KG)") {
		t.Errorf("FormatProduct minimal: %q", got)
	}
	if strings.Contains(got, "[]") {
		t.Errorf("FormatProduct should skip empty manufacturer: %q", got)
	}
}

func TestFormatCart_Empty(t *testing.T) {
	got := FormatCart(willys.Cart{})
	if got != "Cart is empty." {
		t.Errorf("FormatCart empty = %q, want %q", got, "Cart is empty.")
	}
}

func TestFormatCart(t *testing.T) {
	cart := willys.Cart{
		TotalUnitCount: 2,
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
	for _, want := range []string{"Cart (2 items):", "Mjölk", "x2", "Ägg 15p", "[Garant Eko]", "Total: 85,80 kr"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatCart missing %q in:\n%s", want, got)
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
	for _, want := range []string{"1 orders:", "#3057837654", "2026-03-24", "delivered", "3 291,29 kr"} {
		if !strings.Contains(got, want) {
			t.Errorf("FormatOrderHistory missing %q in:\n%s", want, got)
		}
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
	for _, want := range []string{"#3057837654", "Levererad", "3 291,29 kr", "Mejeri:", "Mjölk Färsk 3%", "x2", "(100010649_ST)"} {
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
			{Title: "Frukt", URL: "/frukt", Children: nil},
		},
	}
	got := FormatCategory(cat, 0)
	if !strings.Contains(got, "Frukt & Grönt (/frukt-gront)") {
		t.Errorf("FormatCategory root: %q", got)
	}
	if !strings.Contains(got, "  Frukt (/frukt)") {
		t.Errorf("FormatCategory child: %q", got)
	}
}

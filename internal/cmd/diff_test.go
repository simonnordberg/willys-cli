package cmd

import (
	"strings"
	"testing"
)

func TestFormatDiff_Added(t *testing.T) {
	old := map[string]diffItem{}
	new := map[string]diffItem{
		"101476110_ST": {Code: "101476110_ST", Name: "A-fil 3%", Manufacturer: "Garant", Volume: "1kg", NewQty: 4},
	}
	got := FormatDiff(old, new, "order 123", "cart")
	if !strings.Contains(got, "Added (1)") {
		t.Errorf("expected Added section:\n%s", got)
	}
	if !strings.Contains(got, "+ 101476110_ST") {
		t.Errorf("expected + prefix:\n%s", got)
	}
}

func TestFormatDiff_Removed(t *testing.T) {
	old := map[string]diffItem{
		"101476110_ST": {Code: "101476110_ST", Name: "A-fil 3%", Manufacturer: "Garant", Volume: "1kg", OldQty: 4},
	}
	new := map[string]diffItem{}
	got := FormatDiff(old, new, "order 123", "cart")
	if !strings.Contains(got, "Removed (1)") {
		t.Errorf("expected Removed section:\n%s", got)
	}
	if !strings.Contains(got, "- 101476110_ST") {
		t.Errorf("expected - prefix:\n%s", got)
	}
}

func TestFormatDiff_Changed(t *testing.T) {
	old := map[string]diffItem{
		"101476110_ST": {Code: "101476110_ST", Name: "A-fil 3%", OldQty: 2},
	}
	new := map[string]diffItem{
		"101476110_ST": {Code: "101476110_ST", Name: "A-fil 3%", NewQty: 4},
	}
	got := FormatDiff(old, new, "order 123", "cart")
	if !strings.Contains(got, "Changed (1)") {
		t.Errorf("expected Changed section:\n%s", got)
	}
	if !strings.Contains(got, "~ 101476110_ST") || !strings.Contains(got, "2 →  4") {
		t.Errorf("expected ~ prefix with qty change:\n%s", got)
	}
}

func TestFormatDiff_NoDifferences(t *testing.T) {
	items := map[string]diffItem{
		"101476110_ST": {Code: "101476110_ST", Name: "A-fil", OldQty: 4, NewQty: 4},
	}
	got := FormatDiff(items, items, "order 123", "order 456")
	if !strings.Contains(got, "No differences") {
		t.Errorf("expected no differences:\n%s", got)
	}
}

func TestFormatDiff_Mixed(t *testing.T) {
	old := map[string]diffItem{
		"AAA_ST": {Code: "AAA_ST", Name: "Removed item", OldQty: 1},
		"BBB_ST": {Code: "BBB_ST", Name: "Changed item", OldQty: 2},
	}
	new := map[string]diffItem{
		"BBB_ST": {Code: "BBB_ST", Name: "Changed item", NewQty: 5},
		"CCC_ST": {Code: "CCC_ST", Name: "Added item", NewQty: 3},
	}
	got := FormatDiff(old, new, "order 100", "order 200")
	if !strings.Contains(got, "Added (1)") || !strings.Contains(got, "Removed (1)") || !strings.Contains(got, "Changed (1)") {
		t.Errorf("expected all three sections:\n%s", got)
	}
}

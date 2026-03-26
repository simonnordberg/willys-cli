package cmd

import (
	"fmt"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/willys"
)

// formatDesc builds "Name [Manufacturer] Volume" for any product-like thing.
func formatDesc(name, manufacturer, volume string) string {
	parts := []string{name}
	if manufacturer != "" {
		parts = append(parts, fmt.Sprintf("[%s]", manufacturer))
	}
	if volume != "" {
		parts = append(parts, volume)
	}
	return strings.Join(parts, " ")
}

// FormatProduct formats a search/browse result line.
// Format: CODE  PRICE  (COMPARE)  Description
func FormatProduct(p willys.Product) string {
	compare := ""
	if p.ComparePrice != "" && p.ComparePriceUnit != "" {
		compare = fmt.Sprintf("%s/%s", p.ComparePrice, p.ComparePriceUnit)
	}
	desc := formatDesc(p.Name, p.Manufacturer, p.DisplayVolume)
	if compare != "" {
		return fmt.Sprintf("  %-16s  %10s  %-14s  %s", p.Code, p.Price, compare, desc)
	}
	return fmt.Sprintf("  %-16s  %10s                  %s", p.Code, p.Price, desc)
}

// FormatCart formats the full cart.
// Format per line: CODE  QTY  PRICE  (COMPARE)  Description
func FormatCart(cart willys.Cart) string {
	if cart.TotalUnitCount == 0 {
		return "Cart is empty."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Cart — %d items\n\n", cart.TotalUnitCount)
	for _, p := range cart.Products {
		compare := ""
		if p.ComparePrice != "" && p.ComparePriceUnit != "" {
			compare = fmt.Sprintf("%s/%s", p.ComparePrice, p.ComparePriceUnit)
		}
		desc := formatDesc(p.Name, p.Manufacturer, p.DisplayVolume)
		if compare != "" {
			fmt.Fprintf(&b, "  %-16s  %2d  %10s  %-14s  %s\n", p.Code, p.PickQuantity, p.TotalPrice, compare, desc)
		} else {
			fmt.Fprintf(&b, "  %-16s  %2d  %10s                  %s\n", p.Code, p.PickQuantity, p.TotalPrice, desc)
		}
	}
	fmt.Fprintf(&b, "\nTotal: %s", cart.TotalPrice)
	return b.String()
}

// FormatOrderHistory formats the order list.
// Format: ORDER_NUMBER  DATE  STATUS  TOTAL
func FormatOrderHistory(orders []willys.OrderSummary) string {
	if len(orders) == 0 {
		return "No orders found."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d orders\n\n", len(orders))
	for _, o := range orders {
		status := o.OrderStatus.Code
		if status == "" {
			status = "unknown"
		}
		price := o.Total
		if price == "" {
			price = o.TotalPrice.FormattedValue
		}
		date := o.DeliveryDate
		if date == "" {
			date = o.OrderDate
		}
		fmt.Fprintf(&b, "  %-12s  %-10s  %-10s  %s\n", o.OrderNumber, date, status, price)
	}
	return strings.TrimRight(b.String(), "\n")
}

// FormatOrderDetail formats a single order with items grouped by category.
// Format: CODE  QTY  PRICE  Description
func FormatOrderDetail(o willys.OrderDetail) string {
	var b strings.Builder
	status := o.StatusDisplay
	if status == "" {
		status = o.OrderStatus.Code
	}
	total := o.TotalPrice.FormattedValue
	if total == "" {
		total = o.NettoCost.FormattedValue
	}
	fmt.Fprintf(&b, "Order %s — %s — %s\n", o.OrderNumber, status, total)
	for category, entries := range o.Products {
		fmt.Fprintf(&b, "\n%s\n", category)
		for _, e := range entries {
			qty := e.PickQuantity
			if qty == 0 {
				qty = e.Quantity
			}
			desc := formatDesc(e.Name, e.Manufacturer, e.DisplayVolume)
			fmt.Fprintf(&b, "  %-16s  %2d  %10s  %s\n", e.Code, qty, e.TotalPrice, desc)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// FormatCategory formats the category tree.
func FormatCategory(cat willys.Category, depth int) string {
	if depth > 2 {
		return ""
	}
	var b strings.Builder
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(&b, "%s%-30s  %s\n", indent, cat.Title, cat.URL)
	for _, child := range cat.Children {
		b.WriteString(FormatCategory(child, depth+1))
	}
	return b.String()
}

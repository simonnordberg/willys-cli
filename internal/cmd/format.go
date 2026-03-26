package cmd

import (
	"fmt"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/willys"
)

func FormatProduct(p willys.Product) string {
	parts := []string{p.Name}
	if p.Manufacturer != "" {
		parts = append(parts, fmt.Sprintf("[%s]", p.Manufacturer))
	}
	if p.DisplayVolume != "" {
		parts = append(parts, p.DisplayVolume)
	}
	parts = append(parts, fmt.Sprintf("— %s", p.Price))
	if p.ComparePrice != "" && p.ComparePriceUnit != "" {
		parts = append(parts, fmt.Sprintf("(%s/%s)", p.ComparePrice, p.ComparePriceUnit))
	}
	parts = append(parts, fmt.Sprintf("(%s)", p.Code))
	return "  " + strings.Join(parts, " ")
}

func FormatCart(cart willys.Cart) string {
	if cart.TotalUnitCount == 0 {
		return "Cart is empty."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Cart (%d items):\n", cart.TotalUnitCount)
	for _, p := range cart.Products {
		parts := []string{fmt.Sprintf("  %s", p.Name)}
		if p.Manufacturer != "" {
			parts = append(parts, fmt.Sprintf("[%s]", p.Manufacturer))
		}
		if p.DisplayVolume != "" {
			parts = append(parts, p.DisplayVolume)
		}
		parts = append(parts, fmt.Sprintf("x%d — %s", p.PickQuantity, p.TotalPrice))
		if p.ComparePrice != "" && p.ComparePriceUnit != "" {
			parts = append(parts, fmt.Sprintf("(%s/%s)", p.ComparePrice, p.ComparePriceUnit))
		}
		parts = append(parts, fmt.Sprintf("(%s)", p.Code))
		fmt.Fprintln(&b, strings.Join(parts, " "))
	}
	fmt.Fprintf(&b, "Total: %s", cart.TotalPrice)
	return b.String()
}

func FormatOrderHistory(orders []willys.OrderSummary) string {
	if len(orders) == 0 {
		return "No orders found."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d orders:\n", len(orders))
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
		fmt.Fprintf(&b, "  #%s  %s  %s  %s\n", o.OrderNumber, date, status, price)
	}
	return strings.TrimRight(b.String(), "\n")
}

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
	fmt.Fprintf(&b, "Order #%s (%s)\n", o.OrderNumber, status)
	if o.DeliveryDate != "" {
		fmt.Fprintf(&b, "Delivery: %s\n", o.DeliveryDate)
	}
	if total != "" {
		fmt.Fprintf(&b, "Total: %s\n", total)
	}
	fmt.Fprintln(&b)
	for category, entries := range o.Products {
		fmt.Fprintf(&b, "  %s:\n", category)
		for _, e := range entries {
			qty := e.PickQuantity
			if qty == 0 {
				qty = e.Quantity
			}
			parts := []string{fmt.Sprintf("    %s", e.Name)}
			if e.Manufacturer != "" {
				parts = append(parts, fmt.Sprintf("[%s]", e.Manufacturer))
			}
			if e.DisplayVolume != "" {
				parts = append(parts, e.DisplayVolume)
			}
			if qty > 1 {
				parts = append(parts, fmt.Sprintf("x%d", qty))
			}
			parts = append(parts, fmt.Sprintf("— %s", e.TotalPrice))
			parts = append(parts, fmt.Sprintf("(%s)", e.Code))
			fmt.Fprintln(&b, strings.Join(parts, " "))
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func FormatCategory(cat willys.Category, depth int) string {
	if depth > 2 {
		return ""
	}
	var b strings.Builder
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(&b, "%s%s (%s)\n", indent, cat.Title, cat.URL)
	for _, child := range cat.Children {
		b.WriteString(FormatCategory(child, depth+1))
	}
	return b.String()
}

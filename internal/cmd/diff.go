package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/willys"
	"github.com/spf13/cobra"
)

type diffItem struct {
	Code         string
	Name         string
	Manufacturer string
	Volume       string
	OldQty       int
	NewQty       int
}

func diffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <order-id> [order-id]",
		Short: "Compare cart vs order, or two orders",
		Long: `Compare current cart against a past order, or two orders against each other.

  willys diff 3057837654              # cart vs order
  willys diff 3057837654 3056473722   # order vs order`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}

			var oldItems, newItems map[string]diffItem
			var oldLabel, newLabel string

			if len(args) == 2 {
				// Order vs order
				oldItems, oldLabel, err = orderToDiffItems(c, args[0])
				if err != nil {
					return err
				}
				newItems, newLabel, err = orderToDiffItems(c, args[1])
				if err != nil {
					return err
				}
			} else {
				// Cart vs order
				oldItems, oldLabel, err = orderToDiffItems(c, args[0])
				if err != nil {
					return err
				}
				newItems, newLabel = cartToDiffItems(c)
			}

			fmt.Println(FormatDiff(oldItems, newItems, oldLabel, newLabel))
			return nil
		},
	}
}

func orderToDiffItems(c *willys.Client, orderNumber string) (map[string]diffItem, string, error) {
	order, err := c.GetOrderDetail(orderNumber)
	if err != nil {
		return nil, "", err
	}
	items := make(map[string]diffItem)
	for _, entries := range order.Products {
		for _, e := range entries {
			qty := e.PickQuantity
			if qty == 0 {
				qty = e.Quantity
			}
			items[e.Code] = diffItem{
				Code:         e.Code,
				Name:         e.Name,
				Manufacturer: e.Manufacturer,
				Volume:       e.DisplayVolume,
				OldQty:       qty,
				NewQty:       qty,
			}
		}
	}
	return items, "order " + order.OrderNumber, nil
}

func cartToDiffItems(c *willys.Client) (map[string]diffItem, string) {
	cart, err := c.GetCart()
	if err != nil {
		return make(map[string]diffItem), "cart"
	}
	items := make(map[string]diffItem)
	for _, p := range cart.Products {
		items[p.Code] = diffItem{
			Code:         p.Code,
			Name:         p.Name,
			Manufacturer: p.Manufacturer,
			Volume:       p.DisplayVolume,
			OldQty:       p.PickQuantity,
			NewQty:       p.PickQuantity,
		}
	}
	return items, "cart"
}

func FormatDiff(oldItems, newItems map[string]diffItem, oldLabel, newLabel string) string {
	var added, removed, changed []diffItem

	for code, n := range newItems {
		if o, ok := oldItems[code]; ok {
			if n.NewQty != o.OldQty {
				changed = append(changed, diffItem{
					Code:         code,
					Name:         n.Name,
					Manufacturer: n.Manufacturer,
					Volume:       n.Volume,
					OldQty:       o.OldQty,
					NewQty:       n.NewQty,
				})
			}
		} else {
			added = append(added, n)
		}
	}
	for code, o := range oldItems {
		if _, ok := newItems[code]; !ok {
			removed = append(removed, o)
		}
	}

	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	sort.Slice(removed, func(i, j int) bool { return removed[i].Name < removed[j].Name })
	sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })

	var b strings.Builder
	fmt.Fprintf(&b, "%s → %s\n", oldLabel, newLabel)

	if len(added) > 0 {
		fmt.Fprintf(&b, "\nAdded (%d)\n", len(added))
		for _, d := range added {
			fmt.Fprintf(&b, "  + %-16s  %2d  %s\n", d.Code, d.NewQty, formatDesc(d.Name, d.Manufacturer, d.Volume))
		}
	}
	if len(removed) > 0 {
		fmt.Fprintf(&b, "\nRemoved (%d)\n", len(removed))
		for _, d := range removed {
			fmt.Fprintf(&b, "  - %-16s  %2d  %s\n", d.Code, d.OldQty, formatDesc(d.Name, d.Manufacturer, d.Volume))
		}
	}
	if len(changed) > 0 {
		fmt.Fprintf(&b, "\nChanged (%d)\n", len(changed))
		for _, d := range changed {
			fmt.Fprintf(&b, "  ~ %-16s  %2d → %2d  %s\n", d.Code, d.OldQty, d.NewQty, formatDesc(d.Name, d.Manufacturer, d.Volume))
		}
	}
	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		fmt.Fprint(&b, "\nNo differences.")
	}

	return strings.TrimRight(b.String(), "\n")
}

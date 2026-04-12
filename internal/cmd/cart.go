package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cartCmd() *cobra.Command {
	cart := &cobra.Command{
		Use:   "cart",
		Short: "Shopping cart operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cartList()
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Show cart contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cartList()
		},
	}

	addCmd := &cobra.Command{
		Use:   "add <product-code>",
		Short: "Add a product to the cart",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qty, _ := cmd.Flags().GetInt("qty")
			c, err := GetClient()
			if err != nil {
				return err
			}
			cart, err := c.AddToCart(args[0], qty)
			if err != nil {
				return err
			}
			fmt.Printf("Added %dx %s\n", qty, args[0])
			fmt.Println(FormatCart(cart))
			return nil
		},
	}
	addCmd.Flags().Int("qty", 1, "quantity to add")

	removeCmd := &cobra.Command{
		Use:   "remove <product-code>",
		Short: "Remove a product from the cart",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}
			cart, err := c.RemoveFromCart(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Removed %s\n", args[0])
			fmt.Println(FormatCart(cart))
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the cart",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}
			if err := c.ClearCart(); err != nil {
				return err
			}
			fmt.Println("Cart cleared.")
			cart, err := c.GetCart()
			if err != nil {
				return err
			}
			fmt.Println(FormatCart(cart))
			return nil
		},
	}

	dealsCmd := &cobra.Command{
		Use:   "deals",
		Short: "Show active deals for cart items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cartDeals()
		},
	}

	cart.AddCommand(listCmd, addCmd, removeCmd, clearCmd, dealsCmd)
	return cart
}

func cartDeals() error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	cart, err := c.GetCart()
	if err != nil {
		return err
	}
	var deals []DealProduct
	for _, cp := range cart.Products {
		p, err := c.GetProduct(cp.Code)
		if err != nil {
			continue
		}
		if p.SavingsAmount == nil || *p.SavingsAmount <= 0 || len(p.PotentialPromotions) == 0 {
			continue
		}
		deals = append(deals, DealProduct{
			Code:         cp.Code,
			Name:         cp.Name,
			Manufacturer: cp.Manufacturer,
			Volume:       cp.DisplayVolume,
			Quantity:     cp.PickQuantity,
			TotalPrice:   cp.TotalPrice,
			Promotion:    p.PotentialPromotions[0],
			Savings:      *p.SavingsAmount,
		})
	}
	fmt.Println(FormatCartDeals(deals))
	return nil
}

func cartList() error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	cart, err := c.GetCart()
	if err != nil {
		return err
	}
	fmt.Println(FormatCart(cart))
	return nil
}

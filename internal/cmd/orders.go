package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ordersCmd() *cobra.Command {
	orders := &cobra.Command{
		Use:   "orders",
		Short: "View order history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ordersList()
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all orders",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ordersList()
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <order-number>",
		Short: "Show order details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}
			order, err := c.GetOrderDetail(args[0])
			if err != nil {
				return err
			}
			fmt.Println(FormatOrderDetail(order))
			return nil
		},
	}

	orders.AddCommand(listCmd, showCmd)
	return orders
}

func ordersList() error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	orders, err := c.GetOrderHistory()
	if err != nil {
		return err
	}
	fmt.Println(FormatOrderHistory(orders))
	return nil
}

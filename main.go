package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/env"
	"github.com/simonnordberg/willys-cli/internal/tui"
	"github.com/simonnordberg/willys-cli/internal/willys"
	"github.com/spf13/cobra"
)

var batch string

func main() {
	root := &cobra.Command{
		Use:   "willys",
		Short: "Willys.se grocery store CLI",
		Long:  "Search products, browse categories, and manage your shopping cart at Willys.se.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			env.Load(".env")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if batch != "" {
				return runBatch(batch)
			}
			// No subcommand — launch TUI.
			c, err := getClient()
			if err != nil {
				return err
			}
			return tui.Run(c)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.Flags().StringVarP(&batch, "batch", "i", "", "CSV file with batch operations")

	root.AddCommand(loginCmd())
	root.AddCommand(logoutCmd())
	root.AddCommand(statusCmd())
	root.AddCommand(searchCmd())
	root.AddCommand(categoriesCmd())
	root.AddCommand(browseCmd())
	root.AddCommand(cartCmd())
	root.AddCommand(ordersCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getCreds() (string, string, error) {
	username := strings.Trim(os.Getenv("WILLYS_USERNAME"), `"`)
	password := strings.Trim(os.Getenv("WILLYS_PASSWORD"), `"`)
	if username == "" || password == "" {
		return "", "", fmt.Errorf("credentials required: set WILLYS_USERNAME and WILLYS_PASSWORD env vars or use a .env file")
	}
	return username, password, nil
}

func getClient() (*willys.Client, error) {
	c := willys.NewClient()

	if c.IsLoggedIn() {
		return c, nil
	}

	username, password, err := getCreds()
	if err != nil {
		return nil, err
	}
	if _, err := c.Login(username, password); err != nil {
		return nil, err
	}
	return c, nil
}

func loginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in and save session",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, password, err := getCreds()
			if err != nil {
				return err
			}
			c := willys.NewClient()
			cust, err := c.Login(username, password)
			if err != nil {
				return err
			}
			fmt.Printf("Logged in as %s %s (%s)\n", cust.FirstName, cust.LastName, cust.Email)
			return nil
		},
	}
}

func logoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear saved session",
		RunE: func(cmd *cobra.Command, args []string) error {
			willys.ClearSession()
			fmt.Println("Session cleared.")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check login status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := willys.NewClient()
			if c.IsLoggedIn() {
				cust, err := c.GetCustomer()
				if err != nil {
					return err
				}
				fmt.Printf("Logged in as %s %s (%s)\n", cust.FirstName, cust.LastName, cust.Email)
			} else {
				fmt.Println("Not logged in.")
			}
			return nil
		},
	}
}

func searchCmd() *cobra.Command {
	var count int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for products",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}


			query := strings.Join(args, " ")
			var collected []willys.Product
			page := 0
			pageSize := min(count, 30)

			for len(collected) < count {
				result, err := c.Search(query, page, pageSize)
				if err != nil {
					return err
				}
				collected = append(collected, result.Results...)
				if page == 0 {
					fmt.Printf("%d results for %q:\n", result.Pagination.TotalNumberOfResults, query)
				}
				if page+1 >= result.Pagination.NumberOfPages {
					break
				}
				page++
			}
			for _, p := range collected[:min(count, len(collected))] {
				fmt.Println(formatProduct(p))
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&count, "count", "n", 10, "number of results")
	return cmd
}

func categoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List product categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}


			tree, err := c.Categories()
			if err != nil {
				return err
			}
			printCategory(tree, 0)
			return nil
		},
	}
}

func browseCmd() *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:   "browse <category-path>",
		Short: "Browse products in a category",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}


			result, err := c.Browse(args[0], page, 10)
			if err != nil {
				return err
			}
			fmt.Printf("%d products:\n", result.Pagination.TotalNumberOfResults)
			for _, p := range result.Results {
				fmt.Println(formatProduct(p))
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	return cmd
}

func cartCmd() *cobra.Command {
	cart := &cobra.Command{
		Use:   "cart",
		Short: "Shopping cart operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}


			cart, err := c.GetCart()
			if err != nil {
				return err
			}
			printCart(cart)
			return nil
		},
	}

	addCmd := &cobra.Command{
		Use:   "add <product-code>",
		Short: "Add a product to the cart",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qty, _ := cmd.Flags().GetInt("qty")
			c, err := getClient()
			if err != nil {
				return err
			}


			cart, err := c.AddToCart(args[0], qty)
			if err != nil {
				return err
			}
			fmt.Printf("Added %dx %s\n", qty, args[0])
			printCart(cart)
			return nil
		},
	}
	addCmd.Flags().Int("qty", 1, "quantity to add")

	removeCmd := &cobra.Command{
		Use:   "remove <product-code>",
		Short: "Remove a product from the cart",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}


			cart, err := c.RemoveFromCart(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Removed %s\n", args[0])
			printCart(cart)
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the cart",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
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
			printCart(cart)
			return nil
		},
	}

	cart.AddCommand(addCmd, removeCmd, clearCmd)
	return cart
}

func ordersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "orders [order-number]",
		Short: "View order history or order details",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				order, err := c.GetOrderDetail(args[0])
				if err != nil {
					return err
				}
				printOrderDetail(order)
				return nil
			}
			orders, err := c.GetOrderHistory()
			if err != nil {
				return err
			}
			printOrderHistory(orders)
			return nil
		},
	}
}

func printOrderHistory(orders []willys.OrderSummary) {
	if len(orders) == 0 {
		fmt.Println("No orders found.")
		return
	}
	fmt.Printf("%d orders:\n", len(orders))
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
		fmt.Printf("  #%s  %s  %s  %s\n", o.OrderNumber, date, status, price)
	}
}

func printOrderDetail(o willys.OrderDetail) {
	status := o.StatusDisplay
	if status == "" {
		status = o.OrderStatus.Code
	}
	total := o.TotalPrice.FormattedValue
	if total == "" {
		total = o.NettoCost.FormattedValue
	}
	fmt.Printf("Order #%s (%s)\n", o.OrderNumber, status)
	if o.DeliveryDate != "" {
		fmt.Printf("Delivery: %s\n", o.DeliveryDate)
	}
	if total != "" {
		fmt.Printf("Total: %s\n", total)
	}
	fmt.Println()
	for category, entries := range o.Products {
		fmt.Printf("  %s:\n", category)
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
			fmt.Println(strings.Join(parts, " "))
		}
	}
}

func formatProduct(p willys.Product) string {
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

func printCart(cart willys.Cart) {
	if cart.TotalUnitCount == 0 {
		fmt.Println("Cart is empty.")
		return
	}
	fmt.Printf("Cart (%d items):\n", cart.TotalUnitCount)
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
		fmt.Println(strings.Join(parts, " "))
	}
	fmt.Printf("Total: %s\n", cart.TotalPrice)
}

func printCategory(cat willys.Category, depth int) {
	if depth > 2 {
		return
	}
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s (%s)\n", indent, cat.Title, cat.URL)
	for _, child := range cat.Children {
		printCategory(child, depth+1)
	}
}

func runBatch(file string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ",")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		op := fields[0]
		args := fields[1:]

		fmt.Printf("> %s %s\n", op, strings.Join(args, " "))
		if err := runOp(c, op, args); err != nil {
			return fmt.Errorf("batch operation %q: %w", line, err)
		}
		fmt.Println()
	}
	return scanner.Err()
}

func runOp(c *willys.Client, op string, args []string) error {
	switch op {
	case "cart":
		cart, err := c.GetCart()
		if err != nil {
			return err
		}
		printCart(cart)
	case "add":
		if len(args) < 1 {
			return fmt.Errorf("add requires a product code")
		}
		qty := 1
		if len(args) > 1 {
			q, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid quantity: %s", args[1])
			}
			qty = q
		}
		cart, err := c.AddToCart(args[0], qty)
		if err != nil {
			return err
		}
		fmt.Printf("Added %dx %s\n", qty, args[0])
		printCart(cart)
	case "remove":
		if len(args) < 1 {
			return fmt.Errorf("remove requires a product code")
		}
		cart, err := c.RemoveFromCart(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Removed %s\n", args[0])
		printCart(cart)
	case "clear":
		if err := c.ClearCart(); err != nil {
			return err
		}
		fmt.Println("Cart cleared.")
		cart, err := c.GetCart()
		if err != nil {
			return err
		}
		printCart(cart)
	case "search":
		if len(args) < 1 {
			return fmt.Errorf("search requires a query")
		}
		count := 10
		if len(args) > 1 {
			c, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid count: %s", args[1])
			}
			count = c
		}
		// Reuse the search query — but we don't have a Client variable named c here
		// since the parameter shadows it. The client is passed as the first arg.
		return runSearch(c, args[0], count)
	case "categories":
		tree, err := c.Categories()
		if err != nil {
			return err
		}
		printCategory(tree, 0)
	case "orders":
		if len(args) > 0 {
			order, err := c.GetOrderDetail(args[0])
			if err != nil {
				return err
			}
			printOrderDetail(order)
		} else {
			orders, err := c.GetOrderHistory()
			if err != nil {
				return err
			}
			printOrderHistory(orders)
		}
	case "browse":
		if len(args) < 1 {
			return fmt.Errorf("browse requires a category path")
		}
		page := 0
		if len(args) > 1 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid page: %s", args[1])
			}
			page = p
		}
		result, err := c.Browse(args[0], page, 10)
		if err != nil {
			return err
		}
		fmt.Printf("%d products:\n", result.Pagination.TotalNumberOfResults)
		for _, p := range result.Results {
			fmt.Println(formatProduct(p))
		}
	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
	return nil
}

func runSearch(c *willys.Client, query string, count int) error {
	var collected []willys.Product
	page := 0
	pageSize := min(count, 30)
	for len(collected) < count {
		result, err := c.Search(query, page, pageSize)
		if err != nil {
			return err
		}
		collected = append(collected, result.Results...)
		if page == 0 {
			fmt.Printf("%d results for %q:\n", result.Pagination.TotalNumberOfResults, query)
		}
		if page+1 >= result.Pagination.NumberOfPages {
			break
		}
		page++
	}
	for _, p := range collected[:min(count, len(collected))] {
		fmt.Println(formatProduct(p))
	}
	return nil
}

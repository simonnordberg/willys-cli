package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/willys"
)

func RunBatch(file string) error {
	c, err := GetClient()
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
		fmt.Println(FormatCart(cart))
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
		fmt.Println(FormatCart(cart))
	case "remove":
		if len(args) < 1 {
			return fmt.Errorf("remove requires a product code")
		}
		cart, err := c.RemoveFromCart(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Removed %s\n", args[0])
		fmt.Println(FormatCart(cart))
	case "clear":
		if err := c.ClearCart(); err != nil {
			return err
		}
		fmt.Println("Cart cleared.")
		cart, err := c.GetCart()
		if err != nil {
			return err
		}
		fmt.Println(FormatCart(cart))
	case "search":
		if len(args) < 1 {
			return fmt.Errorf("search requires a query")
		}
		count := 10
		if len(args) > 1 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid count: %s", args[1])
			}
			count = n
		}
		return runSearch(c, args[0], count)
	case "categories":
		tree, err := c.Categories()
		if err != nil {
			return err
		}
		fmt.Print(FormatCategory(tree, 0))
	case "orders":
		if len(args) > 0 {
			order, err := c.GetOrderDetail(args[0])
			if err != nil {
				return err
			}
			fmt.Println(FormatOrderDetail(order))
		} else {
			orders, err := c.GetOrderHistory()
			if err != nil {
				return err
			}
			fmt.Println(FormatOrderHistory(orders))
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
			fmt.Println(FormatProduct(p))
		}
	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
	return nil
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/willys"
	"github.com/spf13/cobra"
)

func searchCmd() *cobra.Command {
	var count int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for products",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}
			query := strings.Join(args, " ")
			return runSearch(c, query, count)
		},
	}
	cmd.Flags().IntVarP(&count, "count", "n", 10, "number of results")
	return cmd
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
		fmt.Println(FormatProduct(p))
	}
	return nil
}

func categoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List product categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := GetClient()
			if err != nil {
				return err
			}
			tree, err := c.Categories()
			if err != nil {
				return err
			}
			fmt.Print(FormatCategory(tree, 0))
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
			c, err := GetClient()
			if err != nil {
				return err
			}
			result, err := c.Browse(args[0], page, 10)
			if err != nil {
				return err
			}
			fmt.Printf("%d products:\n", result.Pagination.TotalNumberOfResults)
			for _, p := range result.Results {
				fmt.Println(FormatProduct(p))
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	return cmd
}

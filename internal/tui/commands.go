package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/simonnordberg/willys-cli/internal/willys"
)

// cartQty returns the current quantity of a product in the cart.
func cartQty(cart willys.Cart, code string) int {
	for _, p := range cart.Products {
		if p.Code == code {
			return p.PickQuantity
		}
	}
	return 0
}

// Messages

type customerMsg struct {
	customer willys.Customer
	err      error
}

type searchResultMsg struct {
	result willys.SearchResult
	query  string
	err    error
}

type categoriesMsg struct {
	root willys.Category
	err  error
}

type browseResultMsg struct {
	result willys.SearchResult
	err    error
}

type cartMsg struct {
	cart willys.Cart
	err  error
}

type addedToCartMsg struct {
	cart willys.Cart
	code string
	err  error
}

type cartClearedMsg struct {
	cart willys.Cart
	err  error
}

type orderHistoryMsg struct {
	orders []willys.OrderSummary
	err    error
}

type orderDetailMsg struct {
	order willys.OrderDetail
	err   error
}

type reorderMsg struct {
	cart willys.Cart
	err  error
}

type statusMsg string

// Commands

func fetchCustomerCmd(c *willys.Client) tea.Cmd {
	return func() tea.Msg {
		cust, err := c.GetCustomer()
		return customerMsg{cust, err}
	}
}

func searchProductsCmd(c *willys.Client, query string, page int) tea.Cmd {
	return func() tea.Msg {
		result, err := c.Search(query, page, 20)
		return searchResultMsg{result, query, err}
	}
}

func fetchCategoriesCmd(c *willys.Client) tea.Cmd {
	return func() tea.Msg {
		root, err := c.Categories()
		return categoriesMsg{root, err}
	}
}

func browseProductsCmd(c *willys.Client, path string, page int) tea.Cmd {
	return func() tea.Msg {
		result, err := c.Browse(path, page, 20)
		return browseResultMsg{result, err}
	}
}

func fetchCartCmd(c *willys.Client) tea.Cmd {
	return func() tea.Msg {
		cart, err := c.GetCart()
		return cartMsg{cart, err}
	}
}

func addToCartCmd(c *willys.Client, code string, qty int) tea.Cmd {
	return func() tea.Msg {
		cart, err := c.AddToCart(code, qty)
		return addedToCartMsg{cart, code, err}
	}
}

func removeFromCartCmd(c *willys.Client, code string) tea.Cmd {
	return func() tea.Msg {
		cart, err := c.RemoveFromCart(code)
		return addedToCartMsg{cart, code, err}
	}
}

func clearCartCmd(c *willys.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.ClearCart(); err != nil {
			return cartClearedMsg{willys.Cart{}, err}
		}
		cart, err := c.GetCart()
		return cartClearedMsg{cart, err}
	}
}

func fetchOrderHistoryCmd(c *willys.Client) tea.Cmd {
	return func() tea.Msg {
		orders, err := c.GetOrderHistory()
		return orderHistoryMsg{orders, err}
	}
}

func fetchOrderDetailCmd(c *willys.Client, orderNumber string) tea.Cmd {
	return func() tea.Msg {
		order, err := c.GetOrderDetail(orderNumber)
		return orderDetailMsg{order, err}
	}
}

func reorderCmd(c *willys.Client, products map[string][]willys.OrderEntry) tea.Cmd {
	return func() tea.Msg {
		var cart willys.Cart
		var err error
		for _, items := range products {
			for _, e := range items {
				if e.Code == "" {
					continue
				}
				qty := max(e.PickQuantity, e.Quantity, 1)
				cart, err = c.AddToCart(e.Code, qty)
				if err != nil {
					return reorderMsg{cart, err}
				}
			}
		}
		return reorderMsg{cart, nil}
	}
}

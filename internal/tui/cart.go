package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonnordberg/willys-cli/internal/willys"
)

// cartItem wraps a CartProduct for the bubbles list.
type cartItem struct {
	product willys.CartProduct
}

func (i cartItem) Title() string {
	name := fmt.Sprintf("%s x%d", i.product.Name, i.product.PickQuantity)
	if i.product.Manufacturer != "" {
		name += " [" + i.product.Manufacturer + "]"
	}
	return name
}

func (i cartItem) Description() string {
	desc := priceStyle.Render(i.product.TotalPrice)
	if i.product.ComparePrice != "" && i.product.ComparePriceUnit != "" {
		desc += mutedStyle.Render(fmt.Sprintf(" (%s/%s)", i.product.ComparePrice, i.product.ComparePriceUnit))
	}
	desc += mutedStyle.Render(fmt.Sprintf(" %s", i.product.Code))
	return desc
}

func (i cartItem) FilterValue() string { return i.product.Name }

type cartModel struct {
	items   list.Model
	spinner spinner.Model
	loading bool
}

func newCartModel() cartModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("205"))

	l := list.New(nil, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return cartModel{
		items:   l,
		spinner: sp,
	}
}

func (m cartModel) Update(msg tea.Msg, client *willys.Client, cart willys.Cart) (cartModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Remove):
			if item, ok := m.items.SelectedItem().(cartItem); ok {
				m.loading = true
				return m, removeFromCartCmd(client, item.product.Code)
			}
		case key.Matches(msg, keys.Inc):
			if item, ok := m.items.SelectedItem().(cartItem); ok {
				m.loading = true
				return m, addToCartCmd(client, item.product.Code, item.product.PickQuantity+1)
			}
		case key.Matches(msg, keys.Dec):
			if item, ok := m.items.SelectedItem().(cartItem); ok {
				if item.product.PickQuantity > 1 {
					m.loading = true
					return m, addToCartCmd(client, item.product.Code, item.product.PickQuantity-1)
				}
			}
		case key.Matches(msg, keys.Clear):
			if cart.TotalUnitCount > 0 {
				m.loading = true
				return m, clearCartCmd(client)
			}
		}
		var cmd tea.Cmd
		m.items, cmd = m.items.Update(msg)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *cartModel) setCartItems(cart willys.Cart) {
	m.loading = false
	items := make([]list.Item, len(cart.Products))
	for i, p := range cart.Products {
		items[i] = cartItem{p}
	}
	m.items.SetItems(items)
}

func (m cartModel) View(width, height int, cart willys.Cart) string {
	if m.loading {
		return m.spinner.View() + " Updating cart...\n"
	}

	if cart.TotalUnitCount == 0 {
		return mutedStyle.Render("Cart is empty") + "\n"
	}

	m.items.SetSize(width, height-3)
	footer := "\n" + titleStyle.Render(fmt.Sprintf("Total: %s (%d items)", cart.TotalPrice, cart.TotalUnitCount))
	return m.items.View() + footer
}

func (m *cartModel) setSize(w, h int) {
	m.items.SetSize(w, h-3)
}

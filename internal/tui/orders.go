package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonnordberg/willys-cli/internal/willys"
)

// orderSummaryItem wraps an OrderSummary for the bubbles list.
type orderSummaryItem struct {
	order willys.OrderSummary
}

func (i orderSummaryItem) Title() string {
	title := fmt.Sprintf("Order #%s", i.order.OrderNumber)
	if i.order.DeliveryDate != "" {
		title += " — " + i.order.DeliveryDate
	} else if i.order.OrderDate != "" {
		title += " — " + i.order.OrderDate
	}
	return title
}

func (i orderSummaryItem) Description() string {
	status := i.order.OrderStatus.Code
	if status == "" {
		status = "Unknown"
	}
	price := i.order.Total
	if price == "" {
		price = i.order.TotalPrice.FormattedValue
	}
	parts := []string{status}
	if price != "" {
		parts = append(parts, priceStyle.Render(price))
	}
	return strings.Join(parts, " — ")
}

func (i orderSummaryItem) FilterValue() string { return i.order.OrderNumber }

// orderEntryItem wraps an OrderEntry for the bubbles list.
type orderEntryItem struct {
	entry willys.OrderEntry
}

func (i orderEntryItem) Title() string {
	name := i.entry.Name
	if i.entry.Manufacturer != "" {
		name += " [" + i.entry.Manufacturer + "]"
	}
	if i.entry.DisplayVolume != "" {
		name += " " + i.entry.DisplayVolume
	}
	qty := i.entry.PickQuantity
	if qty == 0 {
		qty = i.entry.Quantity
	}
	if qty > 1 {
		name += fmt.Sprintf(" x%d", qty)
	}
	return name
}

func (i orderEntryItem) Description() string {
	desc := priceStyle.Render(i.entry.TotalPrice)
	desc += mutedStyle.Render(fmt.Sprintf(" %s", i.entry.Code))
	return desc
}

func (i orderEntryItem) FilterValue() string { return i.entry.Name }

const (
	ordersListMode = iota
	orderDetailMode
)

type ordersModel struct {
	orderList    list.Model
	detailList   list.Model
	spinner      spinner.Model
	mode         int
	loading      bool
	loaded       bool
	currentOrder willys.OrderDetail
}

func newOrdersModel() ordersModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("205"))

	ol := list.New(nil, delegate, 0, 0)
	ol.SetShowTitle(false)
	ol.SetShowStatusBar(false)
	ol.SetShowHelp(false)
	ol.SetFilteringEnabled(false)
	ol.DisableQuitKeybindings()

	dl := list.New(nil, delegate, 0, 0)
	dl.SetShowTitle(false)
	dl.SetShowStatusBar(false)
	dl.SetShowHelp(false)
	dl.SetFilteringEnabled(false)
	dl.DisableQuitKeybindings()

	return ordersModel{
		orderList:  ol,
		detailList: dl,
		spinner:    sp,
	}
}

func (m ordersModel) Update(msg tea.Msg, client *willys.Client, cart willys.Cart) (ordersModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if m.mode == ordersListMode {
				if item, ok := m.orderList.SelectedItem().(orderSummaryItem); ok {
					m.mode = orderDetailMode
					m.loading = true
					return m, fetchOrderDetailCmd(client, item.order.OrderNumber)
				}
			}
		case key.Matches(msg, keys.Back):
			if m.mode == orderDetailMode {
				m.mode = ordersListMode
				return m, nil
			}
		case key.Matches(msg, keys.Reorder):
			if m.mode == orderDetailMode && len(m.currentOrder.Products) > 0 {
				m.loading = true
				return m, reorderCmd(client, m.currentOrder.Products)
			}
		}

		if m.mode == ordersListMode {
			var cmd tea.Cmd
			m.orderList, cmd = m.orderList.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.detailList, cmd = m.detailList.Update(msg)
			cmds = append(cmds, cmd)
		}

	case orderHistoryMsg:
		m.loading = false
		if msg.err != nil {
			return m, nil
		}
		m.loaded = true
		items := make([]list.Item, len(msg.orders))
		for i, o := range msg.orders {
			items[i] = orderSummaryItem{o}
		}
		m.orderList.SetItems(items)

	case orderDetailMsg:
		m.loading = false
		if msg.err != nil {
			return m, nil
		}
		m.currentOrder = msg.order
		var items []list.Item
		for _, entries := range msg.order.Products {
			for _, e := range entries {
				items = append(items, orderEntryItem{e})
			}
		}
		m.detailList.SetItems(items)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ordersModel) View(width, height int) string {
	if m.loading && !m.loaded {
		return m.spinner.View() + " Loading orders...\n"
	}

	if m.mode == orderDetailMode {
		if m.loading {
			return m.spinner.View() + " Loading order details...\n"
		}
		o := m.currentOrder
		header := titleStyle.Render(fmt.Sprintf("Order #%s", o.OrderNumber))
		if o.DeliveryDate != "" {
			header += mutedStyle.Render(" — " + o.DeliveryDate)
		}
		status := o.StatusDisplay
		if status == "" {
			status = o.OrderStatus.Code
		}
		header += "\n" + mutedStyle.Render(status)
		total := o.TotalPrice.FormattedValue
		if total == "" {
			total = o.NettoCost.FormattedValue
		}
		if total != "" {
			header += "  " + priceStyle.Render(total)
		}
		header += "\n\n"

		m.detailList.SetSize(width, height-5)
		return header + m.detailList.View()
	}

	if !m.loaded {
		return mutedStyle.Render("No orders") + "\n"
	}

	m.orderList.SetSize(width, height-1)
	return m.orderList.View()
}

func (m *ordersModel) setSize(w, h int) {
	m.orderList.SetSize(w, h-1)
	m.detailList.SetSize(w, h-5)
}

func (m ordersModel) helpKeys() string {
	if m.mode == orderDetailMode {
		return "↑↓: navigate  a: reorder  esc: back"
	}
	return "↑↓: navigate  enter: view details"
}

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonnordberg/willys-cli/internal/willys"
)

const (
	tabSearch = iota
	tabBrowse
	tabCart
	tabCount
)

var tabNames = []string{"Search", "Browse", "Cart"}

// Model is the root TUI model.
type Model struct {
	client   *willys.Client
	customer willys.Customer
	cart     willys.Cart

	activeTab int
	width     int
	height    int

	search   searchModel
	browse   browseModel
	cartView cartModel

	status  string
	err     error
	loading bool
}

func newModel(client *willys.Client) Model {
	return Model{
		client:   client,
		search:   newSearchModel(),
		browse:   newBrowseModel(),
		cartView: newCartModel(),
	}
}

// Run launches the TUI.
func Run(client *willys.Client) error {
	p := tea.NewProgram(newModel(client), tea.WithAltScreen())
	_, err := p.Run()
	// Reset terminal state — bubbletea doesn't always restore cleanly in tmux.
	fmt.Print("\033[?25h\033[0m\033[?1049l")
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchCustomerCmd(m.client),
		fetchCartCmd(m.client),
		fetchCategoriesCmd(m.client),
		m.search.spinner.Tick,
		m.browse.spinner.Tick,
		m.cartView.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentHeight := m.height - 4 // status + tabs + help
		m.search.setSize(m.width, contentHeight)
		m.browse.setSize(m.width, contentHeight)
		m.cartView.setSize(m.width, contentHeight)
		return m, nil

	case tea.KeyMsg:
		// Global keys handled before tab delegation.
		switch {
		case key.Matches(msg, keys.Quit):
			// Don't quit if search input is focused — let it handle the key.
			if m.activeTab == tabSearch && m.search.focused {
				break
			}
			return m, tea.Quit
		case key.Matches(msg, keys.Tab):
			m.activeTab = (m.activeTab + 1) % tabCount
			if m.activeTab == tabCart {
				cmds = append(cmds, fetchCartCmd(m.client))
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, keys.ShiftTab):
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
			if m.activeTab == tabCart {
				cmds = append(cmds, fetchCartCmd(m.client))
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, keys.Search):
			m.activeTab = tabSearch
			m.search.focused = true
			m.search.input.Focus()
			return m, nil
		}

	// Global message handling (shared state).
	case customerMsg:
		if msg.err == nil {
			m.customer = msg.customer
		}

	case cartMsg:
		if msg.err == nil {
			m.cart = msg.cart
			m.cartView.setCartItems(msg.cart)
		}

	case categoriesMsg:
		m.browse, _ = m.browse.Update(msg, m.client, m.cart)

	case browseResultMsg:
		m.browse, _ = m.browse.Update(msg, m.client, m.cart)

	case addedToCartMsg:
		m.search.loading = false
		m.browse.loading = false
		m.cartView.loading = false
		if msg.err == nil {
			m.cart = msg.cart
			m.cartView.setCartItems(msg.cart)
			m.status = fmt.Sprintf("Added %s to cart", msg.code)
		} else {
			m.status = errorStyle.Render(msg.err.Error())
		}

	case cartClearedMsg:
		m.cartView.loading = false
		if msg.err == nil {
			m.cart = msg.cart
			m.cartView.setCartItems(msg.cart)
			m.status = "Cart cleared"
		}

	case statusMsg:
		m.status = string(msg)
	}

	// Delegate to active tab.
	switch m.activeTab {
	case tabSearch:
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg, m.client, m.cart)
		cmds = append(cmds, cmd)
	case tabBrowse:
		var cmd tea.Cmd
		m.browse, cmd = m.browse.Update(msg, m.client, m.cart)
		cmds = append(cmds, cmd)
	case tabCart:
		var cmd tea.Cmd
		m.cartView, cmd = m.cartView.Update(msg, m.client, m.cart)
		cmds = append(cmds, cmd)
	}

	// Forward spinner ticks to all tabs.
	if _, ok := msg.(spinner.TickMsg); ok {
		var cmd tea.Cmd
		if m.activeTab != tabSearch {
			m.search.spinner, cmd = m.search.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.activeTab != tabBrowse {
			m.browse.spinner, cmd = m.browse.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.activeTab != tabCart {
			m.cartView.spinner, cmd = m.cartView.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Status bar.
	userName := m.customer.FirstName
	if userName == "" {
		userName = "Not logged in"
	}
	cartInfo := fmt.Sprintf("Cart: %d items", m.cart.TotalUnitCount)
	statusLeft := "● " + userName
	statusRight := cartInfo
	gap := m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight) - 2
	if gap < 1 {
		gap = 1
	}
	b.WriteString(statusBarStyle.Width(m.width).Render(statusLeft + strings.Repeat(" ", gap) + statusRight))
	b.WriteString("\n")

	// Tab bar.
	var tabs []string
	for i, name := range tabNames {
		if i == m.activeTab {
			tabs = append(tabs, tabActiveStyle.Render("[ "+name+" ]"))
		} else {
			tabs = append(tabs, tabInactiveStyle.Render("  "+name+"  "))
		}
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	b.WriteString("\n\n")

	// Content.
	contentHeight := m.height - 5
	switch m.activeTab {
	case tabSearch:
		b.WriteString(m.search.View(m.width, contentHeight))
	case tabBrowse:
		b.WriteString(m.browse.View(m.width, contentHeight))
	case tabCart:
		b.WriteString(m.cartView.View(m.width, contentHeight, m.cart))
	}

	// Help bar at bottom.
	help := helpStyle.Render(m.helpText())
	// Pad to fill remaining space, then add help.
	currentHeight := lipgloss.Height(b.String())
	remaining := m.height - currentHeight - 1
	if remaining > 0 {
		b.WriteString(strings.Repeat("\n", remaining))
	}
	b.WriteString(help)

	return b.String()
}

func (m Model) helpText() string {
	base := "tab: switch  q: quit"
	if m.status != "" {
		base = m.status + "  |  " + base
	}
	switch m.activeTab {
	case tabSearch:
		if m.search.focused {
			return "enter: search  |  " + base
		}
		return "↑↓: navigate  a: add to cart  /: search  |  " + base
	case tabBrowse:
		return "↑↓: navigate" + m.browse.helpKeys() + "  |  " + base
	case tabCart:
		return "↑↓: navigate  +/-: qty  d: remove  x: clear  |  " + base
	}
	return base
}

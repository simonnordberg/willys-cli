package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonnordberg/willys-cli/internal/willys"
)

// productItem wraps a Product for the bubbles list.
type productItem struct {
	product willys.Product
}

func (i productItem) Title() string {
	name := i.product.Name
	if i.product.Manufacturer != "" {
		name += " [" + i.product.Manufacturer + "]"
	}
	if i.product.DisplayVolume != "" {
		name += " " + i.product.DisplayVolume
	}
	return name
}

func (i productItem) Description() string {
	desc := priceStyle.Render(i.product.Price)
	if i.product.ComparePrice != "" && i.product.ComparePriceUnit != "" {
		desc += mutedStyle.Render(fmt.Sprintf(" (%s/%s)", i.product.ComparePrice, i.product.ComparePriceUnit))
	}
	desc += mutedStyle.Render(fmt.Sprintf(" %s", i.product.Code))
	return desc
}

func (i productItem) FilterValue() string { return i.product.Name }

type searchModel struct {
	input      textinput.Model
	results    list.Model
	spinner    spinner.Model
	query      string
	totalHits  int
	hasResults bool
	loading    bool
	focused    bool // whether the text input is focused
}

func newSearchModel() searchModel {
	ti := textinput.New()
	ti.Placeholder = "Search products..."
	ti.Focus()
	ti.CharLimit = 100

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

	return searchModel{
		input:   ti,
		results: l,
		spinner: sp,
		focused: true,
	}
}

func (m searchModel) Update(msg tea.Msg, client *willys.Client, cart willys.Cart) (searchModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focused {
			switch {
			case key.Matches(msg, keys.Enter):
				q := m.input.Value()
				if q != "" {
					m.query = q
					m.loading = true
					m.focused = false
					m.input.Blur()
					return m, searchProductsCmd(client, q, 0)
				}
			case key.Matches(msg, keys.Back):
				if m.input.Value() == "" {
					return m, nil
				}
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, keys.Search):
				m.focused = true
				m.input.Focus()
				return m, nil
			case key.Matches(msg, keys.Add):
				if item, ok := m.results.SelectedItem().(productItem); ok {
					m.loading = true
					qty := cartQty(cart, item.product.Code) + 1
					return m, addToCartCmd(client, item.product.Code, qty)
				}
			}
			var cmd tea.Cmd
			m.results, cmd = m.results.Update(msg)
			cmds = append(cmds, cmd)
		}

	case searchResultMsg:
		m.loading = false
		if msg.err != nil {
			return m, nil
		}
		m.hasResults = true
		m.totalHits = msg.result.Pagination.TotalNumberOfResults
		items := make([]list.Item, len(msg.result.Results))
		for i, p := range msg.result.Results {
			items[i] = productItem{p}
		}
		m.results.SetItems(items)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m searchModel) View(width, height int) string {
	var s string

	inputView := m.input.View()
	s += inputView + "\n\n"

	if m.loading {
		s += m.spinner.View() + " Searching...\n"
	} else if m.hasResults {
		s += mutedStyle.Render(fmt.Sprintf("%d results for %q", m.totalHits, m.query)) + "\n\n"
		m.results.SetSize(width, height-5)
		s += m.results.View()
	} else {
		s += mutedStyle.Render("Type a search query and press Enter") + "\n"
	}

	return s
}

func (m *searchModel) setSize(w, h int) {
	m.input.Width = w - 4
	m.results.SetSize(w, h-5)
}


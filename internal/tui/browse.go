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

// categoryItem wraps a Category for the bubbles list.
type categoryItem struct {
	category willys.Category
}

func (i categoryItem) Title() string       { return i.category.Title }
func (i categoryItem) Description() string { return mutedStyle.Render(i.category.URL) }
func (i categoryItem) FilterValue() string { return i.category.Title }

const (
	browseCategories = iota
	browseProducts
)

type browseModel struct {
	catList    list.Model
	prodList   list.Model
	spinner    spinner.Model
	categories []willys.Category // current level children
	breadcrumb []struct {
		title    string
		children []willys.Category
	}
	mode    int
	loading bool
	loaded  bool
}

func newBrowseModel() browseModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("205"))

	cl := list.New(nil, delegate, 0, 0)
	cl.SetShowTitle(false)
	cl.SetShowStatusBar(false)
	cl.SetShowHelp(false)
	cl.SetFilteringEnabled(false)
	cl.DisableQuitKeybindings()

	pl := list.New(nil, delegate, 0, 0)
	pl.SetShowTitle(false)
	pl.SetShowStatusBar(false)
	pl.SetShowHelp(false)
	pl.SetFilteringEnabled(false)
	pl.DisableQuitKeybindings()

	return browseModel{
		catList:  cl,
		prodList: pl,
		spinner:  sp,
	}
}

func (m browseModel) Update(msg tea.Msg, client *willys.Client, cart willys.Cart) (browseModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if m.mode == browseCategories {
				if item, ok := m.catList.SelectedItem().(categoryItem); ok {
					cat := item.category
					if len(cat.Children) > 0 {
						// Push current level and drill down.
						m.breadcrumb = append(m.breadcrumb, struct {
							title    string
							children []willys.Category
						}{cat.Title, m.categories})
						m.categories = cat.Children
						m.setCategoryItems()
						return m, nil
					}
					// Leaf category — browse products.
					m.mode = browseProducts
					m.loading = true
					return m, browseProductsCmd(client, cat.URL, 0)
				}
			}
		case key.Matches(msg, keys.Back):
			if m.mode == browseProducts {
				m.mode = browseCategories
				return m, nil
			}
			if len(m.breadcrumb) > 0 {
				last := m.breadcrumb[len(m.breadcrumb)-1]
				m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-1]
				m.categories = last.children
				m.setCategoryItems()
				return m, nil
			}
		case key.Matches(msg, keys.Add):
			if m.mode == browseProducts {
				if item, ok := m.prodList.SelectedItem().(productItem); ok {
					m.loading = true
					qty := cartQty(cart, item.product.Code) + 1
					return m, addToCartCmd(client, item.product.Code, qty)
				}
			}
		}

		if m.mode == browseCategories {
			var cmd tea.Cmd
			m.catList, cmd = m.catList.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.prodList, cmd = m.prodList.Update(msg)
			cmds = append(cmds, cmd)
		}

	case categoriesMsg:
		m.loading = false
		if msg.err != nil {
			return m, nil
		}
		m.loaded = true
		m.categories = msg.root.Children
		m.setCategoryItems()

	case browseResultMsg:
		m.loading = false
		if msg.err != nil {
			return m, nil
		}
		items := make([]list.Item, len(msg.result.Results))
		for i, p := range msg.result.Results {
			items[i] = productItem{p}
		}
		m.prodList.SetItems(items)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *browseModel) setCategoryItems() {
	items := make([]list.Item, len(m.categories))
	for i, cat := range m.categories {
		items[i] = categoryItem{cat}
	}
	m.catList.SetItems(items)
}

func (m browseModel) View(width, height int) string {
	if m.loading && !m.loaded {
		return m.spinner.View() + " Loading categories...\n"
	}

	var header string
	if len(m.breadcrumb) > 0 {
		path := ""
		for _, b := range m.breadcrumb {
			if path != "" {
				path += " > "
			}
			path += b.title
		}
		header = mutedStyle.Render(path+" > ") + "\n\n"
	}

	if m.mode == browseProducts {
		if m.loading {
			return header + m.spinner.View() + " Loading products...\n"
		}
		m.prodList.SetSize(width, height-3)
		return header + m.prodList.View()
	}

	m.catList.SetSize(width, height-3)
	return header + m.catList.View()
}

func (m *browseModel) setSize(w, h int) {
	m.catList.SetSize(w, h-3)
	m.prodList.SetSize(w, h-3)
}

func (m browseModel) helpKeys() string {
	if m.mode == browseProducts {
		return fmt.Sprintf("  %s select  %s add to cart  %s back", keys.Enter.Help().Key, keys.Add.Help().Key, keys.Back.Help().Key)
	}
	if len(m.breadcrumb) > 0 {
		return fmt.Sprintf("  %s select  %s back", keys.Enter.Help().Key, keys.Back.Help().Key)
	}
	return fmt.Sprintf("  %s select", keys.Enter.Help().Key)
}

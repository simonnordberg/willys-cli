# willys-cli

A CLI and interactive TUI for shopping at [Willys.se](https://www.willys.se) (Swedish grocery store).

Search products, browse categories, manage your cart, and review order history — from the terminal.

## Install

```bash
go install github.com/simonnordberg/willys-cli@latest
```

Or build from source:

```bash
git clone https://github.com/simonnordberg/willys-cli.git
cd willys-cli
make build
```

## Setup

Credentials via environment variables or a `.env` file:

```bash
export WILLYS_USERNAME=YYYYMMDDNNNN
export WILLYS_PASSWORD=yourpassword
```

```bash
# Or a .env file in the current directory
echo 'WILLYS_USERNAME=YYYYMMDDNNNN' >> .env
echo 'WILLYS_PASSWORD=yourpassword' >> .env
```

## Usage

### Interactive TUI

```bash
willys-cli
```

Full terminal UI with tabs for search, browse, cart, and orders.

**Controls:** `Tab` switch tabs, `↑↓` navigate, `Enter` select, `a` add to cart, `+/-` adjust quantity, `d` remove, `q` quit.

### CLI commands

```bash
# Search
willys-cli search mjölk
willys-cli search "ekologisk korv" -n 20

# Browse categories
willys-cli categories
willys-cli browse frukt-och-gront/gronsaker

# Cart
willys-cli cart                           # show cart (alias: cart list)
willys-cli cart add 101233933_ST --qty 2
willys-cli cart remove 101233933_ST
willys-cli cart clear

# Orders
willys-cli orders                         # list all (alias: orders list)
willys-cli orders show 3057837654         # order detail

# Session
willys-cli login
willys-cli status
willys-cli logout

# Batch from CSV
willys-cli -i shopping-list.csv
```

### Batch CSV format

```csv
# Weekly shopping
add, 101233933_ST, 2
add, 101205823_ST, 1
search, mjölk
cart
```

Lines starting with `#` are ignored.

### Output format

All commands use a consistent columnar format with product codes first for easy parsing:

```
willys-cli search mjölk
  101205891_ST    13,50 kr  13,50 kr/l      Mjölk 3% [Garant] 1l
  100010649_ST    21,90 kr  14,60 kr/l      Mjölk Färsk 3% [Falköpings] 1,5l

willys-cli cart
  101476110_ST   4    63,60 kr  15,90 kr/kg    A-fil 3% [Garant] 1kg
  100010649_ST   4    87,60 kr  14,60 kr/l     Mjölk Färsk 3% [Falköpings] 1,5l

willys-cli orders show 3057837654
Order 3057837654 — Levererad — 3 291,29 kr

Mejeri, ost & ägg
  101476110_ST   1    15,90 kr  A-fil 3% [Garant] 1kg
  100010649_ST   2    43,80 kr  Mjölk Färsk 3% [Falköpings] 1,5l
```

## AI agent integration

The CLI is designed to be driven by AI agents (e.g. Claude Code). The consistent columnar output with product codes first makes it easy to parse programmatically, and the batch CSV format allows agents to build shopping lists and execute them in one shot.

Example workflow with Claude Code:

1. Agent plans a weekly meal plan based on family preferences
2. Searches for products: `willys-cli search "kokosmjölk"` — parses codes from output
3. Reviews order history: `willys-cli orders show <id>` — reuses known product codes
4. Builds a CSV batch file with all items
5. Executes: `willys-cli -i shopping-list.csv`
6. Verifies: `willys-cli cart` — confirms totals

The `/shop` slash command in the [companion repo](https://github.com/simonnordberg/willys-shopping) automates the full flow — meal planning, product search, cart management, and recipe generation.

## Session

Sessions are persisted to `~/.config/willys-cli/session.json`. After the first login, subsequent commands reuse the saved session until it expires.

## Development

```bash
make build    # build binary
make test     # run tests
make lint     # golangci-lint
```

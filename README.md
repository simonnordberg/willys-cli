# willys-cli

A CLI and interactive TUI for shopping at [Willys.se](https://www.willys.se) (Swedish grocery store).

Search products, browse categories, and manage your shopping cart — from the terminal.

## Install

```bash
go install github.com/simonnordberg/willys-cli@latest
```

Or build from source:

```bash
git clone https://github.com/simonnordberg/willys-cli.git
cd willys-cli
go build -o willys-cli .
```

## Setup

Credentials (Swedish personnummer + password) via environment variables or a `.env` file:

```bash
# Environment variables
export WILLYS_USERNAME=YYYYMMDDNNNN
export WILLYS_PASSWORD=yourpassword

# Or a .env file in the current directory
echo 'WILLYS_USERNAME=YYYYMMDDNNNN' >> .env
echo 'WILLYS_PASSWORD=yourpassword' >> .env
```

## Usage

### Interactive mode (TUI)

```bash
willys-cli
```

Launches a full terminal UI with three tabs:

- **Search** — find products by name
- **Browse** — navigate the category tree
- **Cart** — view and manage your shopping cart

**Controls:** `Tab` to switch tabs, `↑↓` to navigate, `Enter` to select, `a` to add to cart, `+/-` to adjust quantity, `d` to remove, `q` to quit.

### CLI mode

For scripting and automation:

```bash
# Sök efter produkter
willys-cli search bananer
willys-cli search "ekologisk korv" --count 20

# Bläddra i kategorier
willys-cli categories
willys-cli browse frukt-och-gront/gronsaker

# Varukorg
willys-cli cart
willys-cli cart add 101233933_ST --qty 2
willys-cli cart remove 101233933_ST
willys-cli cart clear

# Session
willys-cli login
willys-cli status
willys-cli logout

# Batchoperationer från CSV
willys-cli --batch handlingslista.csv
```

### Batch CSV format

```csv
# Veckans inköp
add,101233933_ST,2
add,101205823_ST,1
remove,101233933_ST
cart
```

Lines starting with `#` are ignored.

## Session

Sessions are persisted to `~/.config/willys-cli/session.json`. After the first login, subsequent commands reuse the saved session until it expires.

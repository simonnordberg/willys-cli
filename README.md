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

Credentials (Swedish personnummer + password) can be provided via:

```bash
# Environment variables
export WILLYS_USERNAME=YYYYMMDDNNNN
export WILLYS_PASSWORD=yourpassword

# Or a .env file in the current directory
echo 'WILLYS_USERNAME=YYYYMMDDNNNN' >> .env
echo 'WILLYS_PASSWORD=yourpassword' >> .env

# Or flags
willys-cli --username YYYYMMDDNNNN --password yourpassword <command>
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
# Search
willys-cli search mjölk
willys-cli search "ekologisk mjölk" --count 20

# Browse categories
willys-cli categories
willys-cli browse frukt-och-gront/frukt

# Cart
willys-cli cart
willys-cli cart add 101233933_ST --qty 2
willys-cli cart remove 101233933_ST
willys-cli cart clear

# Session
willys-cli login
willys-cli status
willys-cli logout

# Batch operations from CSV
willys-cli --batch shopping-list.csv
```

### Batch CSV format

```csv
add,101233933_ST,2
add,101205823_ST,1
remove,101233933_ST
cart
clear
```

Lines starting with `#` are ignored.

## Session

Sessions are persisted to `~/.config/willys-cli/session.json`. After the first login, subsequent commands reuse the saved session until it expires.

# Steel Tables Viewer

A terminal-based viewer for Australian structural steel section properties (AS/NZS standards).

![Steel Tables](steel-tables.svg)

## Features

- **Interactive navigation** — Arrow keys to scroll, page through columns
- **Responsive display** — Adapts to terminal size, fixed headers while scrolling
- **Multiple section types** — UB, UC, WB, WC, PFC, RHS, SHS, CHS, EA, UA tables

## Installation

### From source

```bash
git clone https://github.com/daveh07/steel_tables.git
cd steel_tables
go build -o steel_tables ./cmd/steel_tables
```

### Run directly

```bash
go run ./cmd/steel_tables
```

## Usage

### Interactive mode

```bash
./steel_tables
```

Use the menu to select a table, then navigate with:
- **← →** Page through columns
- **↑ ↓** Scroll rows
- **PgUp/PgDn** Jump pages of rows
- **m** Return to menu
- **q** Quit

### Command-line mode

```bash
./steel_tables UB350
./steel_tables pfc300
```

## Project Structure

```
steel_tables/
├── cmd/
│   └── steel_tables/
│       └── main.go           # Entry point
├── internal/
│   ├── models/
│   │   └── steel.go          # SteelProperty struct
│   ├── columns/
│   │   └── columns.go        # Column definitions & formatters
│   ├── ui/
│   │   ├── colors.go         # Color constants
│   │   ├── terminal_unix.go  # Unix terminal handling
│   │   ├── terminal_windows.go
│   │   ├── header.go         # Header & footer drawing
│   │   ├── menu.go           # Welcome screen
│   │   └── table.go          # Table row rendering
│   └── viewer/
│       └── viewer.go         # Interactive table display
├── data/
│   └── *.json                # Steel property data files
├── go.mod
└── README.md
```

## Data Sources

Steel section properties are based on Australian Standard AS/NZS 3679.1 and manufacturer data.

## License

MIT License - see [LICENSE](LICENSE)

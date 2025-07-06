# Steel Tables Viewer

Interactive terminal-based viewer for steel section properties.

## Quick Start

### Windows
1. Download `steel_tables.exe`, `run_steel_tables.bat`, and the `data` folder
2. Place them all in the same directory
3. **Double-click `run_steel_tables.bat`** to run the program

### macOS
1. Download the appropriate binary for your Mac:
   - `steel_tables_macos_amd64` for Intel Macs
   - `steel_tables_macos_arm64` for Apple Silicon Macs (M1, M2, etc.)
2. Rename the binary to `steel_tables` (remove the platform suffix)
3. Download `run_steel_tables.command` and the `data` folder
4. Place them all in the same directory
5. **Double-click `run_steel_tables.command`** to run the program
   - If macOS blocks it, right-click → "Open" → "Open" to allow it

### Linux
#### Option 1: Desktop File (recommended)
1. Download `steel_tables_linux`, `steel_tables.desktop`, and the `data` folder
2. Rename `steel_tables_linux` to `steel_tables`
3. Place them all in the same directory
4. **Double-click `steel_tables.desktop`** to run the program

#### Option 2: Shell Script
1. Download `steel_tables_linux`, `run_steel_tables.sh`, and the `data` folder
2. Rename `steel_tables_linux` to `steel_tables`
3. Place them all in the same directory
4. **Double-click `run_steel_tables.sh`** (if your file manager supports it)
5. Or run from terminal: `./run_steel_tables.sh`

## Manual Installation

If you prefer to run from the command line:

1. Download the appropriate binary for your platform
2. Download the `data` folder
3. Place them in the same directory
4. Open a terminal in that directory
5. Run:
   - Windows: `steel_tables.exe`
   - macOS/Linux: `./steel_tables`

## Features

- Interactive navigation through steel section tables
- Support for multiple section types (UB, UC, CHS, RHS, etc.)
- Color-coded display with alternating row colors
- Search and filter capabilities
- Comprehensive property data including:
  - Dimensions (depth, width, thickness)
  - Section properties (area, moment of inertia, etc.)
  - Material properties

## Controls

- **↑/↓**: Navigate through lists
- **Enter**: Select item
- **q**: Quit or go back
- **Esc**: Go back to previous screen

## System Requirements

- Windows 10 or later
- macOS 10.12 or later
- Linux with glibc 2.17 or later

## Troubleshooting

### "Permission denied" errors (macOS/Linux)
Make sure the binary is executable:
```bash
chmod +x steel_tables
```

### macOS "App can't be opened" error
Right-click the file → "Open" → "Open" to bypass Gatekeeper.

### Missing data
Ensure the `data` folder is in the same directory as the executable.

## Building from Source

If you have Go installed:

1. Clone or download the source code
2. Run: `go build -o steel_tables .`
3. Or use the build script: `./build_all.sh`

---

For more information or to report issues, visit the project repository.

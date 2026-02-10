# Steel Tables Viewer - Desktop Installation

## ✅ Installation Complete!

The Steel Tables Viewer has been successfully installed and is ready to use.

## How to Use

### Method 1: Desktop Launcher (Recommended)
- Double-click the **Steel-Tables-Viewer.desktop** file on your desktop
- This will open a terminal and start the interactive application
- Type table names like "WC400", "UB350", "PFC300" etc.
- Use arrow keys or < > to navigate between pages
- Press 'm' to return to menu, 'q' to quit

### Method 2: Terminal (Interactive Mode)
```bash
cd to 'Working directory' 
./steel_tables_viewer
```

### Method 3: Terminal (Command Line Mode)
```bash
cd to 'Working directory' 
./steel_tables_viewer WC400
./steel_tables_viewer UB350
./steel_tables_viewer PFC300
```

### Method 4: Launcher Script
```bash
/home/david/Desktop/Programming/Go/steel_tables/steel-tables-launcher.sh
```

## Available Steel Tables
- CHS350, EA300, EA350, PFC300, PFC350
- RHS350, RHS450, SHS350, SHS450
- UA300, UA350, UB300, UB350, UC300, UC350
- WB300, WB350, WC300, WC400

## Navigation Controls (Interactive Mode)
- **Type table name**: Direct access (e.g., "WC400", "UB350")
- **f [filter]**: Filter steel sections (e.g., "f PFC+UB+UC")
- **Arrow Keys** or **< >**: Navigate between column pages in table view
- **m**: Return to main menu from table view
- **q**: Quit application
- **Enter**: Confirm selection

## Filter Examples
- `f PFC` - Shows all PFC sections (380x100PFC, 300x90PFC, etc.)
- `f UB+UC` - Shows all UB and UC sections 
- `f PFC+UB+UC` - Shows all PFC, UB, and UC sections
- `f 200` - Shows all sections containing "200" in the name

## Files Installed
- `steel_tables_viewer` - Main executable
- `Steel-Tables-Viewer.desktop` - Desktop launcher
- `steel-tables-launcher.sh` - Alternative launcher script
- `data/` - Steel property data files

## Notes
- The application works with Australian steel standards
- All measurements are in metric units (mm, kg/m, MPa, etc.)
- Terminal colors are optimized for dark terminals

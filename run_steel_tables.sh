#!/bin/bash
# Linux shell script to run steel_tables in terminal
# Place this file in the same directory as steel_tables (Linux binary)
# Make executable with: chmod +x run_steel_tables.sh

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the script directory
cd "$SCRIPT_DIR"

echo "Starting Steel Tables Viewer..."
echo ""

# Check if the executable exists
if [ ! -f "./steel_tables" ]; then
    echo "ERROR: steel_tables executable not found in this directory!"
    echo "Please make sure steel_tables (Linux binary) is in the same folder as this script."
    echo ""
    read -p "Press Enter to close this window..."
    exit 1
fi

# Make sure the binary is executable
chmod +x ./steel_tables

# Run the steel tables program
./steel_tables

# Keep the terminal open after the program exits
echo ""
echo "Program finished. Press Enter to close this window..."
read

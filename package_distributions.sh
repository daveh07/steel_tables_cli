#!/bin/bash
# Package distributions for easy sharing
# This script creates zip files for each platform

echo "Creating distribution packages..."
echo ""

cd builds

# Create Windows package
echo "Packaging Windows distribution..."
zip -r steel_tables_windows.zip steel_tables.exe run_steel_tables.bat data/
echo "✓ Created steel_tables_windows.zip"

# Create macOS Intel package
echo "Packaging macOS Intel distribution..."
cp steel_tables_macos_amd64 steel_tables
zip -r steel_tables_macos_intel.zip steel_tables run_steel_tables.command data/
rm steel_tables
echo "✓ Created steel_tables_macos_intel.zip"

# Create macOS Apple Silicon package
echo "Packaging macOS Apple Silicon distribution..."
cp steel_tables_macos_arm64 steel_tables
zip -r steel_tables_macos_arm64.zip steel_tables run_steel_tables.command data/
rm steel_tables
echo "✓ Created steel_tables_macos_arm64.zip"

# Create Linux package
echo "Packaging Linux distribution..."
cp steel_tables_linux steel_tables
zip -r steel_tables_linux.zip steel_tables run_steel_tables.sh steel_tables.desktop data/
rm steel_tables
echo "✓ Created steel_tables_linux.zip"

echo ""
echo "Distribution packages created:"
echo "  steel_tables_windows.zip"
echo "  steel_tables_macos_intel.zip"
echo "  steel_tables_macos_arm64.zip"
echo "  steel_tables_linux.zip"
echo ""
echo "Each package contains:"
echo "  - The platform-specific binary"
echo "  - Platform-specific launcher script"
echo "  - Complete data directory"
echo "  - Ready to run by double-clicking the launcher"

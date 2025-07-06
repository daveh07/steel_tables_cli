#!/bin/bash
# Build script to create binaries for all major platforms
# Run this script to build steel_tables for Windows, macOS, and Linux

echo "Building Steel Tables for all platforms..."
echo ""

# Create builds directory if it doesn't exist
mkdir -p builds

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o builds/steel_tables.exe .
if [ $? -eq 0 ]; then
    echo "✓ Windows build successful: builds/steel_tables.exe"
    # Copy the batch file to the builds directory
    cp run_steel_tables.bat builds/
    echo "✓ Copied run_steel_tables.bat to builds/"
else
    echo "✗ Windows build failed"
fi
echo ""

# Build for macOS (64-bit Intel)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o builds/steel_tables_macos_amd64 .
if [ $? -eq 0 ]; then
    echo "✓ macOS (Intel) build successful: builds/steel_tables_macos_amd64"
else
    echo "✗ macOS (Intel) build failed"
fi
echo ""

# Build for macOS (Apple Silicon)
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o builds/steel_tables_macos_arm64 .
if [ $? -eq 0 ]; then
    echo "✓ macOS (Apple Silicon) build successful: builds/steel_tables_macos_arm64"
else
    echo "✗ macOS (Apple Silicon) build failed"
fi

# Copy the command file to builds directory
cp run_steel_tables.command builds/
chmod +x builds/run_steel_tables.command
echo "✓ Copied run_steel_tables.command to builds/"
echo ""

# Build for Linux (64-bit)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o builds/steel_tables_linux .
if [ $? -eq 0 ]; then
    echo "✓ Linux build successful: builds/steel_tables_linux"
    # Copy the shell script and desktop file to builds directory
    cp run_steel_tables.sh builds/
    cp steel_tables.desktop builds/
    chmod +x builds/run_steel_tables.sh
    chmod +x builds/steel_tables.desktop
    echo "✓ Copied run_steel_tables.sh and steel_tables.desktop to builds/"
else
    echo "✗ Linux build failed"
fi
echo ""

# Copy data directory to builds
echo "Copying data directory to builds..."
cp -r data builds/
echo "✓ Data directory copied to builds/"
echo ""

echo "Build process complete!"
echo ""
echo "Distribution packages:"
echo "  Windows: builds/steel_tables.exe + builds/run_steel_tables.bat + builds/data/"
echo "  macOS:   builds/steel_tables_macos_* + builds/run_steel_tables.command + builds/data/"
echo "  Linux:   builds/steel_tables_linux + builds/run_steel_tables.sh + builds/steel_tables.desktop + builds/data/"

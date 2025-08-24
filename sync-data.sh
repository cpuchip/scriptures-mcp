#!/bin/bash

# Script to sync scripture data from the bcbooks/scriptures-json repository
# This keeps our data up to date with the source repository

set -e

echo "Syncing scripture data from bcbooks/scriptures-json repository..."

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Clone the scriptures-json repository
echo "Cloning scriptures-json repository..."
git clone https://github.com/bcbooks/scriptures-json.git

# Copy the data files to our data directory
echo "Copying scripture data files..."
cd scriptures-json

# Get the directory of this script (should be the repository root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATA_DIR="$SCRIPT_DIR/data"

# Create data directory if it doesn't exist
mkdir -p "$DATA_DIR"

# Copy the required files
cp "book-of-mormon.json" "$DATA_DIR/"
cp "doctrine-and-covenants.json" "$DATA_DIR/"
cp "pearl-of-great-price.json" "$DATA_DIR/"
cp "old-testament.json" "$DATA_DIR/"
cp "new-testament.json" "$DATA_DIR/"

echo "Scripture data synchronized successfully!"
echo "Data files updated in: $DATA_DIR"

# Show file sizes
echo ""
echo "Updated files:"
ls -lh "$DATA_DIR"/*.json

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo ""
echo "Data sync complete. You may want to commit these changes."
echo "To see what changed, run: git diff --stat data/"
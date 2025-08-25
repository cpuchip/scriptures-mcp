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

echo "Copying scripture data files..."
cd scriptures-json

# Get the directory of this script (repository root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATA_DIR="$SCRIPT_DIR/internal/scripture/data"

# Create data directory if it doesn't exist
mkdir -p "$DATA_DIR"

# Copy the required files
cp "book-of-mormon.json" "$DATA_DIR/"
cp "doctrine-and-covenants.json" "$DATA_DIR/"
cp "pearl-of-great-price.json" "$DATA_DIR/"
cp "old-testament.json" "$DATA_DIR/"
cp "new-testament.json" "$DATA_DIR/"

# Create compressed archive (zip) for embedding
echo "Creating compressed archive scriptures.zip..."
(
	cd "$DATA_DIR" && \
	if command -v zip >/dev/null 2>&1; then
		zip -9 -q scriptures.zip *.json || echo "Warning: zip failed"
		# If zip succeeded and file exists & non-empty, remove original JSONs
		if [ -s scriptures.zip ]; then
			echo "Removing original JSON files (kept inside scriptures.zip)";
			rm -f *.json
		else
			echo "Warning: scriptures.zip not created properly; keeping JSON files" >&2
		fi
	else
		echo "Warning: 'zip' command not found; scriptures.zip not created. Install zip to enable compressed embedding." >&2
	fi
)

echo "Scripture data synchronized successfully!"
echo "Data files updated in: $DATA_DIR"

# Show file sizes
echo ""
echo "Updated files:"
ls -lh "$DATA_DIR"/*.json 2>/dev/null || true
if [ -f "$DATA_DIR/scriptures.zip" ]; then
	ls -lh "$DATA_DIR/scriptures.zip"
fi

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo ""
echo "Data sync complete. You may want to commit these changes."
echo "To see what changed, run: git diff --stat internal/scripture/data/"
package scripture

import "embed"

// Embedded compressed scripture data (ZIP archive containing JSON files).
// Run ./sync-data.sh (or .\sync-data.ps1) to refresh data/scriptures.zip before building.
//go:embed data/scriptures.zip
var embeddedData embed.FS

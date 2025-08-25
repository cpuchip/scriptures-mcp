package scripture

import "embed"

// Embedded scripture JSON data.
// NOTE: Ensure the JSON files are located in this package folder under data/ before building.
//go:embed data/*.json
var embeddedData embed.FS

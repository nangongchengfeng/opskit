// Package embed manages embedded binary assets
// This file ensures the package can compile even when assets are missing
package embed

// Note: The actual go:embed directives are in manager.go
// This file exists to prevent compilation errors when assets directory is empty

import "embed"

// Dummy variable to ensure the embed package is imported
var _ embed.FS

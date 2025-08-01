package web

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrManifestJSON = errors.New("failed to marshal manifest to JSON")

// Manifest Display Modes
const (
	ManDisplayStandalone = "standalone"
	ManDisplayFullscreen = "fullscreen"
	ManDisplayMinimalUi  = "minimal-ui"
	ManDisplayBrowser    = "browser"
)

// Manifest represents the web app manifest.
type Manifest struct {
	Name            *string   `json:"name"`
	ShortName       *string   `json:"short_name,omitempty"`
	StartURL        *string   `json:"start_url,omitempty"`
	Display         *string   `json:"display,omitempty"`
	BackgroundColor *string   `json:"background_color,omitempty"`
	ThemeColor      *string   `json:"theme_color,omitempty"`
	Icons           []*string `json:"icons,omitempty"`
}

// Build builds the manifest and renders it to JSON format.
func (m *Manifest) Build() (j []byte, err error) {
	// Render the manifest as JSON
	j, err = json.MarshalIndent(m, "", "  ")
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrManifestJSON, err)
	}
	return
}

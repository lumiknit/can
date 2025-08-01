package public

import (
	"embed"
)

//go:embed files
var PublicFS embed.FS

const PublicRoot = "files"

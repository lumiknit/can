package static

import (
	"embed"
	"io/fs"
	"path"
	"strings"
)

type Item struct {
	Path        string
	ContentType string
	Content     []byte
}

//go:embed files
var f embed.FS

func pathToContentType(path string) string {
	dotIdx := strings.LastIndex(path, ".")
	ext := ""
	if dotIdx != -1 {
		ext = path[dotIdx+1:]
	}
	switch ext {
	case "html", "css":
		return "text/" + ext + "; charset=utf-8"
	case "js":
		return "application/javascript; charset=utf-8"
	case "json":
		return "application/json; charset=utf-8"
	case "png", "gif":
		return "image/" + ext
	case "jpg", "jpeg":
		return "image/jpeg"
	case "svg":
		return "image/svg+xml"
	}
	return "application/octet-stream"
}

type Handler func(file *Item)

func eachDir(dirEntries []fs.DirEntry, basePath string, fn Handler) {
	for _, entry := range dirEntries {
		p := path.Join(basePath, entry.Name())
		if entry.IsDir() {
			subEntries, err := f.ReadDir(p)
			if err != nil {
				panic("failed to read embedded static files: " + err.Error())
			}
			eachDir(subEntries, p, fn)
		} else {
			content, err := fs.ReadFile(f, p)
			if err != nil {
				panic("failed to read embedded static file: " + err.Error())
			}
			fn(&Item{
				Path:        strings.TrimPrefix(p, "files/"),
				ContentType: pathToContentType(entry.Name()),
				Content:     content,
			})
		}
	}
}

// Traverse all files in the filesystem
func Each(fn Handler) {
	dirEntries, err := f.ReadDir("files")
	if err != nil {
		panic("failed to read embedded static files: " + err.Error())
	}
	eachDir(dirEntries, "files", fn)
}

package web

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
)

type StaticFile struct {
	Href        string
	ContentType string
	Content     []byte
}

func pathToContentType(p string) string {
	ext := path.Ext(p)
	if len(ext) > 0 {
		ext = ext[1:] // Remove the leading dot
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

type embedFSTraverser struct {
	efs      embed.FS
	rootPath string
	files    []*StaticFile
}

func (t *embedFSTraverser) embedFSToStaticFilesDir(
	relPath string,
	dirEntries []fs.DirEntry,
) error {
	for _, entry := range dirEntries {
		rp := path.Join(relPath, entry.Name())
		fp := path.Join(t.rootPath, rp)
		if entry.IsDir() {
			entries, err := t.efs.ReadDir(fp)
			if err != nil {
				return fmt.Errorf("failed to read subdir '%s': %w", fp, err)
			}
			if err := t.embedFSToStaticFilesDir(rp, entries); err != nil {
				return err
			}
		} else {
			content, err := t.efs.ReadFile(fp)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", fp, err)
			}
			t.files = append(t.files, &StaticFile{
				Href:        StaticPath(rp),
				ContentType: pathToContentType(entry.Name()),
				Content:     content,
			})
		}
	}
	return nil
}

// EmbedFSToStaticFiles converts an embedded filesystem to a slice of StaticFile.
// It reads all files and directories from the embedded filesystem and returns
func EmbedFSToStaticFiles(
	efs embed.FS,
	rootPath string,
) ([]*StaticFile, error) {
	dirEntries, err := efs.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	t := &embedFSTraverser{
		efs:      efs,
		rootPath: rootPath,
		files:    []*StaticFile{},
	}

	err = t.embedFSToStaticFilesDir("", dirEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse embed fs: %w", err)
	}
	return t.files, nil
}

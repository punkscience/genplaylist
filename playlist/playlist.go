package playlist

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Format specifies the playlist output format.
type Format int

const (
	FormatM3U Format = iota
	FormatPLS
)

// Config holds the parameters for playlist generation.
type Config struct {
	Format     Format
	Keywords   []string // empty = no keyword filter
	Extensions []string // empty = default [".mp3", ".flac", ".wav"]
	Title      string   // PLS playlist title (optional, defaults to "playlist")
}

func (c Config) extensions() []string {
	if len(c.Extensions) == 0 {
		return []string{".mp3", ".flac", ".wav"}
	}
	lower := make([]string, len(c.Extensions))
	for i, e := range c.Extensions {
		lower[i] = strings.ToLower(e)
	}
	return lower
}

// Serializer writes a list of file paths to a writer in a playlist format.
type Serializer interface {
	Write(paths []string, w io.Writer) error
}

// NewSerializer returns the Serializer adapter for the given format.
func NewSerializer(f Format, title string) Serializer {
	switch f {
	case FormatPLS:
		return plsSerializer{title: title}
	default:
		return m3uSerializer{}
	}
}

// Generate discovers audio files in fsys and writes a playlist to w.
func Generate(fsys fs.FS, w io.Writer, cfg Config) error {
	paths, err := discoverFiles(fsys, cfg)
	if err != nil {
		return err
	}

	serializer := NewSerializer(cfg.Format, cfg.Title)
	if err := serializer.Write(paths, w); err != nil {
		return fmt.Errorf("writing playlist: %w", err)
	}
	return nil
}

func discoverFiles(fsys fs.FS, cfg Config) ([]string, error) {
	exts := cfg.extensions()
	var paths []string

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !hasAudioExt(path, exts) {
			return nil
		}
		if len(cfg.Keywords) > 0 && !containsKeyword(path, cfg.Keywords) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discovering files: %w", err)
	}
	return paths, nil
}

func hasAudioExt(path string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}

func containsKeyword(path string, keywords []string) bool {
	lowerPath := strings.ToLower(path)
	for _, kw := range keywords {
		if strings.Contains(lowerPath, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

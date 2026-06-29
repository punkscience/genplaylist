package playlist

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type plsSerializer struct {
	title string
}

func (s plsSerializer) Write(paths []string, w io.Writer) error {
	title := s.title
	if title == "" {
		title = "playlist"
	}

	lines := []string{
		"[playlist]",
	}

	for i, p := range paths {
		n := i + 1
		name := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		lines = append(lines,
			fmt.Sprintf("File%d=%s", n, p),
			fmt.Sprintf("Title%d=%s", n, name),
		)
	}
	lines = append(lines, fmt.Sprintf("NumberOfEntries=%d", len(paths)))
	lines = append(lines, "Version=2")

	if _, err := io.WriteString(w, strings.Join(lines, "\n")+"\n"); err != nil {
		return err
	}
	return nil
}

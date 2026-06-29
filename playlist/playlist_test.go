package playlist

import (
	"bytes"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestGenerate_M3U(t *testing.T) {
	fsys := fstest.MapFS{
		"rock/foo.mp3":    {},
		"rock/bar.flac":   {},
		"rock/notes.txt":  {},
		"jazz/baz.wav":    {},
		"jazz/cover.jpg":  {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), buf.String())
	}
	for _, line := range lines {
		if !strings.HasSuffix(line, ".mp3") && !strings.HasSuffix(line, ".flac") && !strings.HasSuffix(line, ".wav") {
			t.Errorf("unexpected line: %q", line)
		}
	}
}

func TestGenerate_PLS(t *testing.T) {
	fsys := fstest.MapFS{
		"track.mp3": {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatPLS, Title: "my mix"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[playlist]") {
		t.Error("missing [playlist] header")
	}
	if !strings.Contains(out, "File1=track.mp3") {
		t.Error("missing File1 entry")
	}
	if !strings.Contains(out, "Title1=track") {
		t.Error("missing Title1 entry")
	}
	if !strings.Contains(out, "NumberOfEntries=1") {
		t.Error("missing NumberOfEntries")
	}
	if !strings.Contains(out, "Version=2") {
		t.Error("missing Version")
	}
}

func TestGenerate_Keywords_CaseInsensitive(t *testing.T) {
	fsys := fstest.MapFS{
		"ROCK/song.MP3":    {},
		"pop/Track.FLAC":   {},
		"jazz/smooth.WAV":  {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U, Keywords: []string{"rock"}})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ROCK/song.MP3") {
		t.Error("keyword 'rock' should match ROCK/song.MP3")
	}
	if strings.Contains(out, "pop") || strings.Contains(out, "jazz") {
		t.Errorf("keyword filter should exclude non-matching files: %q", out)
	}
}

func TestGenerate_NoKeywords_ReturnsAll(t *testing.T) {
	fsys := fstest.MapFS{
		"a.mp3":  {},
		"b.flac": {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 files, got %d", len(lines))
	}
}

func TestGenerate_EmptyDir(t *testing.T) {
	fsys := fstest.MapFS{}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestGenerate_OnlyNonAudio(t *testing.T) {
	fsys := fstest.MapFS{
		"cover.jpg": {},
		"notes.txt": {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestGenerate_CustomExtensions(t *testing.T) {
	fsys := fstest.MapFS{
		"track.ogg":  {},
		"skip.mp3":   {},
	}
	var buf bytes.Buffer

	err := Generate(fsys, &buf, Config{Format: FormatM3U, Extensions: []string{".ogg"}})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "track.ogg") {
		t.Error("should include .ogg file")
	}
	if strings.Contains(out, ".mp3") {
		t.Error("should exclude .mp3 when extensions are custom")
	}
}

func TestGenerate_ErrorWrapped(t *testing.T) {
	fsys := errorFS{}
	err := Generate(fsys, &bytes.Buffer{}, Config{Format: FormatM3U})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "discovering files") {
		t.Errorf("error should wrap context: %v", err)
	}
}

type errorFS struct{}

func (errorFS) Open(name string) (fs.File, error) {
	return nil, errors.New("boom")
}
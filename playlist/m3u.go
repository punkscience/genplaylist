package playlist

import "io"

type m3uSerializer struct{}

func (m3uSerializer) Write(paths []string, w io.Writer) error {
	for _, p := range paths {
		if _, err := io.WriteString(w, p+"\n"); err != nil {
			return err
		}
	}
	return nil
}

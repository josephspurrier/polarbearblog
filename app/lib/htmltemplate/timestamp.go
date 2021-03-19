package htmltemplate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// assetTimePath returns a URL with a timestamp appended.
func assetTimePath(s string) string {
	s = strings.TrimLeft(s, "/")
	abs, err := filepath.Abs(s)
	if err != nil {
		return s
	}

	time, err2 := fileTime(abs)
	if err2 != nil {
		return s
	}

	return fmt.Sprintf("/%v?%v", s, time)
}

// fileTime returns the modification time of the file.
func fileTime(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", fi.ModTime().Unix()), nil
}

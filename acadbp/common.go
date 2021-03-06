package acadbp

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

// BtoYN returns "Y" if b is true, otherwise "N"
func BtoYN(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

// Bto10 returns "1" if b is true, otherwise "0"
func Bto10(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// ExpandGlobPattern returns filenames with the glob patterns expanded from filenames include * or ?
func ExpandGlobPattern(args []string) ([]string, error) {
	ret := []string{}
	for _, arg := range args {
		if strings.ContainsRune(arg, '*') || strings.ContainsRune(arg, '?') {
			g, err := filepath.Glob(arg)
			if err != nil {
				return nil, err
			}
			ret = append(ret, g...)
		} else {
			ret = append(ret, arg)
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("no input file")
	}
	return ret, nil
}

// CreateTempFile creates a new file to temp directory and then writes contents and returns that filepath
func CreateTempFile(name string, contents string, encode string) (string, error) {
	e, err := htmlindex.Get(encode)
	if err != nil {
		return "", err
	}
	temp, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer temp.Close()
	w := transform.NewWriter(temp, e.NewEncoder())
	if _, err := w.Write(bytes.ReplaceAll([]byte(contents), []byte("\n"), []byte("\r\n"))); err != nil {
		return "", err
	}
	return temp.Name(), nil
}

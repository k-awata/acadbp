package acadbp

import (
	"bufio"
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

// StdinToString returns string from stdin
func StdinToString() string {
	s := bufio.NewScanner(os.Stdin)
	var buf bytes.Buffer
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	return buf.String()
}

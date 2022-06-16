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

// BoolToYesNo returns "Y" if b is true, otherwise "N"
func BoolToYesNo(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

// CheckAcCorePath returns error if path of accoreconsole is incorrect
func CheckAcCorePath(accore string) error {
	if !strings.EqualFold(filepath.Base(accore), "accoreconsole.exe") {
		return errors.New("accorepath is incorrect")
	}
	if _, err := os.Stat(accore); err != nil {
		return errors.New("acadbp cannot find accoreconsole binary")
	}
	return nil
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
	return ret, nil
}

// StdinToString returns string from stdin
func StdinToString() string {
	s := bufio.NewScanner(os.Stdin)
	var buf bytes.Buffer
	for s.Scan() {
		buf.WriteString(s.Text() + "\r\n")
	}
	return buf.String()
}

// CreateTempFile creates a new file to temp directory and then writes contents and returns that filepath
func CreateTempFile(name string, contents string, encode string) (string, error) {
	e, err := htmlindex.Get(encode)
	if err != nil {
		return "", err
	}
	str, _, err := transform.String(e.NewEncoder(), contents)
	if err != nil {
		return "", err
	}
	temp, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer temp.Close()
	if _, err := temp.WriteString(str); err != nil {
		return "", err
	}
	return temp.Name(), nil
}

// CreateEmptyFiles creates empty files whose extension is replaced by a specified ext
func CreateEmptyFiles(files []string, ext string) error {
	for _, f := range files {
		if strings.HasPrefix(ext, ".") {
			out := strings.TrimSuffix(f, filepath.Ext(f)) + ext
			if _, err := os.Stat(out); err != nil {
				if _, err := os.Create(out); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

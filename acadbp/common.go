package acadbp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// BoolToYesNo returns "Y" if b is true, otherwise "N"
func BoolToYesNo(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

// ExpandGlobPattern returns filenames with the glob patterns expanded from filenames include * or ?
func ExpandGlobPattern(args []string) []string {
	ret := []string{}
	for _, arg := range args {
		if strings.ContainsRune(arg, '*') || strings.ContainsRune(arg, '?') {
			g, err := filepath.Glob(arg)
			if err != nil {
				continue
			}
			ret = append(ret, g...)
		} else {
			ret = append(ret, arg)
		}
	}
	return ret
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
func CreateTempFile(name string, contents string, sjis bool) (string, error) {
	temp, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer temp.Close()
	// for Japanese
	str := contents
	if sjis {
		str, _, err = transform.String(japanese.ShiftJIS.NewEncoder(), contents)
		if err != nil {
			return "", err
		}
	}
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

// CreateBatContents creates bat file to run accoreconsole and returns that filepath
func CreateBatContents(accore string, scr string, log string, files []string) (string, error) {
	if _, err := os.Stat(accore); err != nil {
		return "", errors.New("acadbp cannot find accoreconsole binary")
	}
	if _, err := os.Stat(scr); err != nil {
		return "", errors.New("acadbp cannot find script file")
	}
	var buf bytes.Buffer
	buf.WriteString("@echo off\r\n")
	buf.WriteString("setlocal\r\n")
	buf.WriteString("set acc=" + accore + "\r\n")
	buf.WriteString("set scr=" + scr + "\r\n")
	buf.WriteString("set log=" + log + "\r\n")
	if len(files) == 0 {
		buf.WriteString(`"%acc%" /s "%scr%" >> "%log%"` + "\r\n")
		return buf.String(), nil
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return "", errors.New("acadbp cannot find drawing " + f)
		}
		buf.WriteString(`"%acc%" /i "` + f + `" /s "%scr%" >> "%log%"` + "\r\n")
	}
	return buf.String(), nil
}

// Runbat executes bat commands
func RunBat(bat string) error {
	temp, err := CreateTempFile("*.bat", bat, true)
	if err != nil {
		return err
	}
	if err := exec.Command("cmd", "/c", temp).Start(); err != nil {
		return err
	}
	fmt.Println("Running accoreconsole...")
	return nil
}

package acadbp

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ReadTemplateDsd returns dsd file contents without [Target], [DWF6Sheet:*] and [MRU *] sections
func ReadTemplateDsd(dsd string, sjis bool) (string, error) {
	file, err := os.Open(dsd)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// for Japanese
	var s *bufio.Scanner
	if sjis {
		s = bufio.NewScanner(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	} else {
		s = bufio.NewScanner(file)
	}
	skip := false
	first := true
	var buf bytes.Buffer
	for s.Scan() {
		str := s.Text()
		// Check if input is dsd format
		if first {
			if str != "[DWF6Version]" {
				return "", errors.New("invalid dsd file")
			}
			first = false
		}
		// Switch if ignoring section
		if strings.HasPrefix(str, "[") {
			skip = str == "[Target]" ||
				strings.HasPrefix(str, "[DWF6Sheet:") ||
				strings.HasPrefix(str, "[MRU ")
		}
		if !skip {
			buf.WriteString(str + "\r\n")
		}
	}
	return buf.String(), nil
}

// CreateDsdTarget returns [Target] section of dsd file with specified options
func CreateDsdTarget(ftype string, multi string) (string, error) {
	// Type=
	no := ""
	ok := false
	if multi == "" {
		no, ok = map[string]string{"dwf": "0", "plotter": "2", "dwfx": "3", "pdf": "5"}[ftype]
	} else {
		no, ok = map[string]string{"dwf": "1", "plotter": "2", "dwfx": "4", "pdf": "6"}[ftype]
	}
	if !ok {
		return "", errors.New("invalid output type")
	}
	// OUT=
	out, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// DWF=
	dwf := filepath.Join(out, "plot.dwf")
	if multi != "" {
		dwf, err = filepath.Abs(multi)
		if err != nil {
			return "", err
		}
	}
	return "[Target]\r\nType=" + no + "\r\nDWF=" + dwf + "\r\nOUT=" + out + string(os.PathSeparator) + "\r\nPWD=\r\n", nil
}

// CreateDsdSheets returns [DWF6Sheet:*] sections of dsd file for each drawing file
func CreateDsdSheets(files []string, sname string, sfile string, layout string) (string, error) {
	if len(files) == 0 {
		return "", errors.New("no input file")
	}
	setup := sname
	if sname != "" && sfile != "" {
		abs, err := filepath.Abs(sfile)
		if err != nil {
			return "", err
		}
		setup = sname + "|" + abs
	}
	var buf bytes.Buffer
	for _, f := range files {
		abs, err := filepath.Abs(f)
		if err != nil {
			return "", err
		}
		buf.WriteString("[DWF6Sheet:" + strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)) + "]\r\n")
		buf.WriteString("DWG=" + abs + "\r\n")
		buf.WriteString("Layout=" + layout + "\r\n")
		buf.WriteString("Setup=" + setup + "\r\n")
		buf.WriteString("OriginalSheetPath=" + abs + "\r\n")
		buf.WriteString("Has Plot Port=0\r\n")
		buf.WriteString("Has3DDWF=0\r\n")
	}
	return buf.String(), nil
}

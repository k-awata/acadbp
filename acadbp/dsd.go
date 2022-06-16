package acadbp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

// ReadTemplateDsd returns dsd file contents without [Target], [DWF6Sheet:*] and [MRU *] sections
func ReadTemplateDsd(dsd string, encode string) (string, error) {
	e, err := htmlindex.Get(encode)
	if err != nil {
		return "", err
	}
	file, err := os.Open(dsd)
	if err != nil {
		return "", err
	}
	defer file.Close()
	str, err := scanTemplateDsd(transform.NewReader(file, e.NewDecoder()))
	if err != nil {
		return "", err
	}
	return str, nil
}

func scanTemplateDsd(r io.Reader) (string, error) {
	s := bufio.NewScanner(r)
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
			buf.WriteString(str + "\n")
		}
	}
	return buf.String(), nil
}

// CreateDsdTarget returns [Target] section of dsd file with specified options
func CreateDsdTarget(ftype string, multi string) (string, error) {
	typeno, err := getDsdTargetType(ftype, multi)
	if err != nil {
		return "", err
	}
	out, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dwf := filepath.Join(out, "plot.dwf")
	if multi != "" {
		dwf, err = filepath.Abs(multi)
		if err != nil {
			return "", err
		}
	}
	return "[Target]\nType=" + typeno + "\nDWF=" + dwf + "\nOUT=" + out + string(os.PathSeparator) + "\nPWD=\n", nil
}

func getDsdTargetType(ftype string, multi string) (string, error) {
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
	return no, nil
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
		buf.WriteString("[DWF6Sheet:" + strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)) + "]\n")
		buf.WriteString("DWG=" + abs + "\n")
		buf.WriteString("Layout=" + layout + "\n")
		buf.WriteString("Setup=" + setup + "\n")
		buf.WriteString("OriginalSheetPath=" + abs + "\n")
		buf.WriteString("Has Plot Port=0\n")
		buf.WriteString("Has3DDWF=0\n")
	}
	return buf.String(), nil
}

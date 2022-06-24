package acadbp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type DsdBuilder struct {
	tmpl   string
	encode encoding.Encoding
	typ    int
	dwf    string
	out    string
	shts   []Sheet
}

// NewDsdBuilder returns a new dsd builder
func NewDsdBuilder() *DsdBuilder {
	return &DsdBuilder{
		tmpl:   "[DWF6Version]\nVer=1\n[DWF6MinorVersion]\nMinorVer=1\n",
		encode: unicode.UTF8,
		typ:    -1,
	}
}

// SetEncoding sets a encoding name to read and write dsd file
func (d *DsdBuilder) SetEncoding(name string) error {
	e, err := htmlindex.Get(name)
	if err != nil {
		return err
	}
	d.encode = e
	return nil
}

// SetTemplate sets a template dsd file
func (d *DsdBuilder) SetTemplate(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := readDsd(transform.NewReader(f, d.encode.NewDecoder()))
	if err != nil {
		return err
	}
	d.tmpl = s
	return nil
}

func readDsd(r io.Reader) (string, error) {
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

// SetOutputFile sets output file type and filename for multi-sheet file
func (d *DsdBuilder) SetOutputFile(ftype string, multifile string) error {
	multi := multifile != ""
	no, err := getDsdType(ftype, multi)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if multi {
		abs, err := filepath.Abs(multifile)
		if err != nil {
			return err
		}
		d.dwf = abs
	} else {
		d.dwf = filepath.Join(wd, "plot.dwf")
	}
	d.typ = no
	d.out = wd + string(os.PathSeparator)
	return nil
}

func getDsdType(t string, m bool) (int, error) {
	no := 0
	ok := false
	if m {
		no, ok = map[string]int{"dwf": 1, "plotter": 2, "dwfx": 4, "pdf": 6}[t]
	} else {
		no, ok = map[string]int{"dwf": 0, "plotter": 2, "dwfx": 3, "pdf": 5}[t]
	}
	if !ok {
		return -1, errors.New("invalid output type")
	}
	return no, nil
}

// SetSheets sets input files to publish and their layout name and page setup
func (d *DsdBuilder) SetSheets(files []string, layout string, sname string, sfile string) error {
	for _, f := range files {
		sh, err := NewSheet(f)
		if err != nil {
			return err
		}
		if err := sh.SetPageSetup(sname, sfile); err != nil {
			return err
		}
		sh.SetLayout(layout)
		d.shts = append(d.shts, *sh)
	}
	return nil
}

// Output returns dsd file contents
func (d *DsdBuilder) Output() string {
	var buf bytes.Buffer
	buf.WriteString(d.tmpl)
	buf.WriteString("[Target]\n")
	buf.WriteString("Type=" + strconv.Itoa(d.typ) + "\n")
	buf.WriteString("DWF=" + d.dwf + "\n")
	buf.WriteString("OUT=" + d.out + "\n")
	buf.WriteString("PWD=\n")
	for _, s := range d.shts {
		buf.WriteString(s.MakeDsdFormat())
	}
	return buf.String()
}

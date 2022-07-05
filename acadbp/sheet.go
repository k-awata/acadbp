package acadbp

import (
	"path/filepath"
	"strings"
)

type Sheet struct {
	path     string
	name     string
	layout   string
	setup    string
	plotport bool
	dwf3d    bool
}

// NewSheet returns a new sheet to publish
func NewSheet(path string) (*Sheet, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &Sheet{
		path:   abs,
		name:   strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		layout: "Model",
	}, nil
}

// SetLayout sets layout name to publish
func (s *Sheet) SetLayout(l string) {
	s.layout = l
}

// SetPageSetup sets a page setup name and a filename that includes the setup
func (s *Sheet) SetPageSetup(name string, file string) error {
	if file == "" {
		s.setup = name
		return nil
	}
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	s.setup = name + "|" + abs
	return nil
}

// MakeDsdFormat returns a sheet setting for dsd file
func (s *Sheet) MakeDsdFormat() string {
	return "[DWF6Sheet:" + s.name + "]\n" +
		"DWG=" + s.path + "\n" +
		"Layout=" + s.layout + "\n" +
		"Setup=" + s.setup + "\n" +
		"OriginalSheetPath=" + s.path + "\n" +
		"Has Plot Port=" + Bto10(s.plotport) + "\n" +
		"Has3DDWF=" + Bto10(s.dwf3d) + "\n"
}

package acadbp

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

package acadbp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

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

// RunBatCommands executes bat commands
func RunBatCommands(cmd string, encode string) error {
	temp, err := CreateTempFile("*.bat", cmd, encode)
	if err != nil {
		return err
	}
	if err := exec.Command("cmd", "/c", temp).Start(); err != nil {
		return err
	}
	fmt.Println("Running accoreconsole...")
	return nil
}

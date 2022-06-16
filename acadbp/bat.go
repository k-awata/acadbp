package acadbp

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CreateBatContents creates bat file to run accoreconsole and returns that filepath
func CreateBatContents(accore string, scr string, log string, files []string) string {
	var buf bytes.Buffer
	buf.WriteString("@echo off\r\n")
	buf.WriteString("setlocal\r\n")
	buf.WriteString("set acc=" + accore + "\r\n")
	buf.WriteString("set scr=" + scr + "\r\n")
	buf.WriteString("set log=" + log + "\r\n")
	if len(files) == 0 {
		buf.WriteString(`"%acc%" /s "%scr%" >> "%log%"` + "\r\n")
		return buf.String()
	}
	for _, f := range files {
		buf.WriteString(`"%acc%" /i "` + f + `" /s "%scr%" >> "%log%"` + "\r\n")
	}
	return buf.String()
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

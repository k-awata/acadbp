package acadbp

import "testing"

func TestCreateBatContents(t *testing.T) {
	type args struct {
		accore string
		scr    string
		log    string
		files  []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"test",
			args{
				"accoreconsole.exe",
				"script.scr",
				"acadbp.log",
				[]string{"Drawing1.dwg", "Drawing2.dwg", "Drawing3.dwg"},
			},
			`@echo off
setlocal
set acc=accoreconsole.exe
set scr=script.scr
set log=acadbp.log
"%acc%" /i "Drawing1.dwg" /s "%scr%" >> "%log%"
"%acc%" /i "Drawing2.dwg" /s "%scr%" >> "%log%"
"%acc%" /i "Drawing3.dwg" /s "%scr%" >> "%log%"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateBatContents(tt.args.accore, tt.args.scr, tt.args.log, tt.args.files); got != tt.want {
				t.Errorf("CreateBatContents() = %v, want %v", got, tt.want)
			}
		})
	}
}

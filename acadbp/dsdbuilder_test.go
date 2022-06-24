package acadbp

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func Test_readDsd(t *testing.T) {
	const input = `[DWF6Version]
Ver=1
[DWF6MinorVersion]
MinorVer=1
[DWF6Sheet:Drawing1-Model]
Layout=Model
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing2-Model]
Layout=Model
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing3-Model]
Layout=Model
Has Plot Port=0
Has3DDWF=0
[Target]
Type=0
PWD=
[PdfOptions]
IncludeHyperlinks=TRUE
CreateBookmarks=TRUE
CaptureFontsInDrawing=TRUE
ConvertTextToGeometry=FALSE
[MRU block template]
MRU=0
[MRU Local]
MRU=1
[MRU Sheet List]
MRU=0
[AutoCAD Block Data]
IncludeBlockInfo=0
BlockTmplFilePath=
`
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Success",
			args{bytes.NewBufferString(input)},
			`[DWF6Version]
Ver=1
[DWF6MinorVersion]
MinorVer=1
[PdfOptions]
IncludeHyperlinks=TRUE
CreateBookmarks=TRUE
CaptureFontsInDrawing=TRUE
ConvertTextToGeometry=FALSE
[AutoCAD Block Data]
IncludeBlockInfo=0
BlockTmplFilePath=
`,
			false,
		},
		{
			"Failure",
			args{bytes.NewBufferString("corrupted")},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readDsd(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDsd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readDsd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDsdType(t *testing.T) {
	type args struct {
		t string
		m bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"DWF s", args{"dwf", false}, 0, false},
		{"DWF m", args{"dwf", true}, 1, false},
		{"Plot s", args{"plotter", false}, 2, false},
		{"Plot m", args{"plotter", true}, 2, false},
		{"DWFX s", args{"dwfx", false}, 3, false},
		{"DWFX m", args{"dwfx", true}, 4, false},
		{"PDF s", args{"pdf", false}, 5, false},
		{"PDF m", args{"pdf", true}, 6, false},
		{"Incorrect", args{"incorrect", false}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDsdType(tt.args.t, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDsdType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDsdType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDsdBuilder_Output(t *testing.T) {
	dsd := NewDsdBuilder()
	dsd.typ = 5
	dsd.dwf = `c:\tmp\plot.dwf`
	dsd.out = `c:\tmp\`
	err := dsd.SetSheets(
		[]string{`c:\tmp\Drawing1.dwg`, `c:\tmp\Drawing2.dwg`, `c:\tmp\Drawing3.dwg`},
		"Layout1",
		"Setup1",
		`c:\tmp\setup.dwg`,
	)
	fmt.Println(dsd)
	if err != nil {
		t.FailNow()
	}
	tests := []struct {
		name string
		d    *DsdBuilder
		want string
	}{
		{
			"test",
			dsd,
			`[DWF6Version]
Ver=1
[DWF6MinorVersion]
MinorVer=1
[Target]
Type=5
DWF=c:\tmp\plot.dwf
OUT=c:\tmp\
PWD=
[DWF6Sheet:Drawing1]
DWG=c:\tmp\Drawing1.dwg
Layout=Layout1
Setup=Setup1|c:\tmp\setup.dwg
OriginalSheetPath=c:\tmp\Drawing1.dwg
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing2]
DWG=c:\tmp\Drawing2.dwg
Layout=Layout1
Setup=Setup1|c:\tmp\setup.dwg
OriginalSheetPath=c:\tmp\Drawing2.dwg
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing3]
DWG=c:\tmp\Drawing3.dwg
Layout=Layout1
Setup=Setup1|c:\tmp\setup.dwg
OriginalSheetPath=c:\tmp\Drawing3.dwg
Has Plot Port=0
Has3DDWF=0
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Output(); got != tt.want {
				t.Errorf("DsdBuilder.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}

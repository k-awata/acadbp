package acadbp

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func Test_scanTemplateDsd(t *testing.T) {
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
		// TODO: Add test cases.
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
			got, err := scanTemplateDsd(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanTemplateDsd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("scanTemplateDsd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDsdTarget(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}
	type args struct {
		ftype string
		multi string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Single PDF",
			args{"pdf", ""},
			`[Target]
Type=5
DWF=` + wd + string(os.PathSeparator) + `plot.dwf
OUT=` + wd + string(os.PathSeparator) + `
PWD=
`,
			false,
		},
		{
			"Multi PDF",
			args{"pdf", "multi.pdf"},
			`[Target]
Type=6
DWF=` + wd + string(os.PathSeparator) + `multi.pdf
OUT=` + wd + string(os.PathSeparator) + `
PWD=
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDsdTarget(tt.args.ftype, tt.args.multi)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDsdTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateDsdTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDsdTargetType(t *testing.T) {
	type args struct {
		ftype string
		multi string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"DWF s", args{"dwf", ""}, "0", false},
		{"DWF m", args{"dwf", "m"}, "1", false},
		{"Plotter s", args{"plotter", ""}, "2", false},
		{"Plotter m", args{"plotter", "m"}, "2", false},
		{"DWFX s", args{"dwfx", ""}, "3", false},
		{"DWFX m", args{"dwfx", "m"}, "4", false},
		{"PDF s", args{"pdf", ""}, "5", false},
		{"PDF m", args{"pdf", "m"}, "6", false},
		{"Incorrect", args{"incorrect", ""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDsdTargetType(tt.args.ftype, tt.args.multi)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDsdTargetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDsdTargetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDsdSheets(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}
	type args struct {
		files  []string
		sname  string
		sfile  string
		layout string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Normal",
			args{
				[]string{"Drawing1.dwg", "Drawing2.dwg", "Drawing3.dwg"},
				"Setup1",
				"setup.dwg",
				"Layout1",
			},
			`[DWF6Sheet:Drawing1]
DWG=` + wd + string(os.PathSeparator) + `Drawing1.dwg
Layout=Layout1
Setup=Setup1|` + wd + string(os.PathSeparator) + `setup.dwg
OriginalSheetPath=` + wd + string(os.PathSeparator) + `Drawing1.dwg
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing2]
DWG=` + wd + string(os.PathSeparator) + `Drawing2.dwg
Layout=Layout1
Setup=Setup1|` + wd + string(os.PathSeparator) + `setup.dwg
OriginalSheetPath=` + wd + string(os.PathSeparator) + `Drawing2.dwg
Has Plot Port=0
Has3DDWF=0
[DWF6Sheet:Drawing3]
DWG=` + wd + string(os.PathSeparator) + `Drawing3.dwg
Layout=Layout1
Setup=Setup1|` + wd + string(os.PathSeparator) + `setup.dwg
OriginalSheetPath=` + wd + string(os.PathSeparator) + `Drawing3.dwg
Has Plot Port=0
Has3DDWF=0
`,
			false,
		},
		{
			"No file",
			args{
				[]string{},
				"Setup1",
				"setup.dwg",
				"Layout1",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDsdSheets(tt.args.files, tt.args.sname, tt.args.sfile, tt.args.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDsdSheets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateDsdSheets() = %v, want %v", got, tt.want)
			}
		})
	}
}

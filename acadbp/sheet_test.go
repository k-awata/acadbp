package acadbp

import "testing"

func TestSheet_MakeDsdFormat(t *testing.T) {
	tests := []struct {
		name string
		s    *Sheet
		want string
	}{
		{
			"test",
			&Sheet{
				path:   `c:\tmp\Drawing1.dwg`,
				name:   "Drawing1",
				layout: "Layout1",
				setup:  `Setup1|c:\tmp\setup.dwg`,
			},
			`[DWF6Sheet:Drawing1]
DWG=c:\tmp\Drawing1.dwg
Layout=Layout1
Setup=Setup1|c:\tmp\setup.dwg
OriginalSheetPath=c:\tmp\Drawing1.dwg
Has Plot Port=0
Has3DDWF=0
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.MakeDsdFormat(); got != tt.want {
				t.Errorf("Sheet.MakeDsdFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

package acadbp

import "testing"

func TestBoolToYesNo(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"true to Yes", args{true}, "Y"},
		{"false to No", args{false}, "N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolToYesNo(tt.args.b); got != tt.want {
				t.Errorf("BoolToYesNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

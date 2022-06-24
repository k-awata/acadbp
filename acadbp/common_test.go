package acadbp

import "testing"

func TestBtoYN(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"true to Yes", args{true}, "Y"},
		{"false to No", args{false}, "N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BtoYN(tt.args.b); got != tt.want {
				t.Errorf("BtoYN() = %v, want %v", got, tt.want)
			}
		})
	}
}

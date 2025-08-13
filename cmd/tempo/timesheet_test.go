package tempo

import (
	"testing"
)

func Test_stringToInt(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Valid Number", args{"65"}, 65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToInt(tt.args.input); got != tt.want {
				t.Errorf("stringToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

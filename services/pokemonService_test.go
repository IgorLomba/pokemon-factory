package services

import (
	"testing"
)

func Test_contains(t *testing.T) {
	type args struct {
		slice []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true if item is in slice",
			args: args{
				slice: []string{"a", "b", "c"},
				item:  "a",
			},
			want: true,
		},
		{
			name: "Should return false if item is not in slice",
			args: args{
				slice: []string{"a", "b", "c"},
				item:  "Gulden Draak",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

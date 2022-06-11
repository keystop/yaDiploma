package luhn

import (
	"testing"
)

func TestCheck(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Тест",
			args: args{
				arr: []int{4, 5, 6, 1, 2, 6, 1, 2, 1, 2, 3, 4, 5, 4, 6, 4},
			},
			want: false,
		},
		{
			name: "Тест",
			args: args{
				arr: []int{4, 5, 6, 1, 2, 6, 1, 2, 1, 2, 3, 4, 5, 4, 6, 7},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.args.arr); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Тест",
			args: args{
				s: "4561261212345464",
			},
			want: false,
		},
		{
			name: "Тест",
			args: args{
				s: "4561261212345467",
			},
			want: true,
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckString(tt.args.s); got != tt.want {
				t.Errorf("CheckString() = %v, want %v", got, tt.want)
			}
		})
	}
}

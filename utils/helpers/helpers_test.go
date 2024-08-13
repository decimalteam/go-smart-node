package helpers

import "testing"

func TestGetBaseDenom(t *testing.T) {
	type args struct {
		chainID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetBaseDenom",
			args: args{
				chainID: "decimal_202020-221213",
			},
			want: "tdel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBaseDenom(tt.args.chainID); got != tt.want {
				t.Errorf("GetBaseDenom() = %v, want %v", got, tt.want)
			}
		})
	}
}

package hash

import "testing"

func TestGetHash(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty input",
			args: args{b: []byte("")},
			want: "da39a3ee",
		},
		{
			name: "simple input",
			args: args{b: []byte("http://ya.ru")},
			want: "9e7174c1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHash(tt.args.b); got != tt.want {
				t.Errorf("GetHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

package kubernetes

import "testing"

func TestTransferDns(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"dns.other",
			args{name: "fe@#@3.12哈"},
			"fe3.12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransferDns(tt.args.name); got != tt.want {
				t.Errorf("TransferDns() = %v, want %v", got, tt.want)
			}
		})
	}
}

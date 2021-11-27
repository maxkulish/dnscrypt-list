package download

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGetRemote(t *testing.T) {
	type args struct {
		remoteURL *url.URL
	}
	tests := []struct {
		name string
		args args
		want *http.Response
	}{
		{
			name: "first test",
			args: args{
				remoteURL:
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRemote(tt.args.remoteURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}

package waiter

import (
	"reflect"
	"testing"
)

func TestGetBytes(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Encode Nil", args{d: nil}, []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBytes(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

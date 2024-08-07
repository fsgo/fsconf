package confext

import (
	"reflect"
	"testing"

	"github.com/fsgo/fsconf"
)

func TestParse(t *testing.T) {
	type args struct {
		confName string
		obj      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				confName: "not_exists.json",
				obj:      map[string]string{},
			},
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				confName: "abc.json",
				obj:      map[string]string{},
			},
			want: map[string]string{
				"A": "bb",
			},
			wantErr: false,
		},
		{
			name: "case 3",
			args: args{
				confName: "db1.toml",
				obj:      map[string]string{},
			},
			want: map[string]string{
				"name":    "abc",
				"charset": "utf-8",
				"Port":    "8080",
			},
			wantErr: false,
		},
		{
			name: "case 4",
			args: args{
				confName: "db1",
				obj:      map[string]string{},
			},
			want: map[string]string{
				"name":    "abc",
				"charset": "utf-8",
				"Port":    "8080",
			},
			wantErr: false,
		},
		{
			name: "case 5",
			args: args{
				confName: "db2", // 存在同名目录的情况
				obj:      map[string]string{},
			},
			want: map[string]string{
				"name": "abc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsconf.Parse(tt.args.confName, &tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.args.obj
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got=%v want=%v", got, tt.want)
			}
		})
	}
}

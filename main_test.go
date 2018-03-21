package main

import (
	"testing"
)

func Test_fixTargetExtension(t *testing.T) {
	type args struct {
		targetPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dir",
			args: args{"dir"},
			want: "dir.zip",
		},
		{
			name: "nested_dir",
			args: args{"/dir/nestedDir"},
			want: "/dir/nestedDir.zip",
		},
		{
			name: "withExtension",
			args: args{"dir.zip"},
			want: "dir.zip",
		},
		{
			name: "file",
			args: args{"file.text"},
			want: "file.text.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixTargetExtension(tt.args.targetPath); got != tt.want {
				t.Errorf("fixTargetExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigsModel_validate(t *testing.T) {
	type fields struct {
		SourcePath string
		TargetPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "file",
			fields:  fields{"./", "./"},
			wantErr: false,
		},
		{
			name:    "not_exist",
			fields:  fields{"./folder", "./"},
			wantErr: true,
		},
		{
			name:    "sourcPath_empty",
			fields:  fields{"", "./folder"},
			wantErr: true,
		},
		{
			name:    "targetPath_empty",
			fields:  fields{"./folder", ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{
				SourcePath: tt.fields.SourcePath,
				TargetPath: tt.fields.TargetPath,
			}
			if err := cfg.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigsModel.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateTargetPath(t *testing.T) {
	type args struct {
		targetPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty",
			args: args{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateTargetPath(tt.args.targetPath)
		})
	}
}

package main

import (
	"testing"
)

func TestFixTargetExtension(t *testing.T) {
	type args struct {
		targetDir string
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
			if got := fixTargetExtension(tt.args.targetDir); got != tt.want {
				t.Errorf("fixTargetExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigsModel_validate(t *testing.T) {
	tests := []struct {
		name       string
		SourcePath string
		TargetDir  string
		wantErr    bool
	}{
		{
			name:       "file",
			SourcePath: "./",
			TargetDir:  "./",
			wantErr:    false,
		},
		{
			name:       "not_exist",
			SourcePath: "./folder",
			TargetDir:  "./",
			wantErr:    true,
		},
		{
			name:       "sourcPath_empty",
			SourcePath: "",
			TargetDir:  "./folder",
			wantErr:    true,
		},
		{
			name:       "targetDir_empty",
			SourcePath: "./folder",
			TargetDir:  "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{
				SourcePath: tt.SourcePath,
				TargetDir:  tt.TargetDir,
			}
			if err := cfg.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigsModel.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	tests := []struct {
		name      string
		targetDir string
	}{
		{
			name:      "empty",
			targetDir: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureDir(tt.targetDir)
		})
	}
}

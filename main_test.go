package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFixDestination(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, base string) (destination string, sourcePath string)
		wantSuffix string
	}{
		{
			name: "appends .zip when destination has no extension",
			setup: func(t *testing.T, base string) (string, string) {
				return filepath.Join(base, "archive"), filepath.Join(base, "src")
			},
			wantSuffix: "archive.zip",
		},
		{
			name: "keeps single .zip when destination already ends in .zip",
			setup: func(t *testing.T, base string) (string, string) {
				return filepath.Join(base, "archive.zip"), filepath.Join(base, "src")
			},
			wantSuffix: "archive.zip",
		},
		{
			name: "joins source basename when destination is an existing directory",
			setup: func(t *testing.T, base string) (string, string) {
				dst := filepath.Join(base, "out")
				if err := os.MkdirAll(dst, 0755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				return dst, filepath.Join(base, "src")
			},
			wantSuffix: filepath.Join("out", "src.zip"),
		},
		{
			name: "preserves source basename extension when joining into a directory",
			setup: func(t *testing.T, base string) (string, string) {
				dst := filepath.Join(base, "out")
				if err := os.MkdirAll(dst, 0755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				return dst, filepath.Join(base, "folder_structure.apk")
			},
			wantSuffix: filepath.Join("out", "folder_structure.apk.zip"),
		},
		{
			name: "cleans the destination path",
			setup: func(t *testing.T, base string) (string, string) {
				return filepath.Join(base, "sub", "..", "archive.zip"), filepath.Join(base, "src")
			},
			wantSuffix: "archive.zip",
		},
		{
			name: "creates missing parent directory",
			setup: func(t *testing.T, base string) (string, string) {
				return filepath.Join(base, "new", "nested", "archive.zip"), filepath.Join(base, "src")
			},
			wantSuffix: filepath.Join("new", "nested", "archive.zip"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base := t.TempDir()
			destination, sourcePath := tc.setup(t, base)

			got, err := fixDestination(destination, sourcePath)
			if err != nil {
				t.Fatalf("fixDestination returned error: %v", err)
			}

			want := filepath.Join(base, tc.wantSuffix)
			if got != want {
				t.Errorf("fixDestination(%q, %q) = %q, want %q", destination, sourcePath, got, want)
			}

			if _, err := os.Stat(filepath.Dir(got)); err != nil {
				t.Errorf("parent directory of result was not created: %v", err)
			}
		})
	}
}

func TestFixDestinationExt(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"archive", "archive.zip"},
		{"archive.zip", "archive.zip"},
		{"folder.apk", "folder.apk.zip"},
		{"path/to/archive", "path/to/archive.zip"},
		{"path/to/archive.zip", "path/to/archive.zip"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := fixDestinationExt(tc.in); got != tc.want {
				t.Errorf("fixDestinationExt(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

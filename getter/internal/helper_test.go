package internal

import (
	"testing"
)

func TestGetAssetMetadata(t *testing.T) {
	tests := []struct {
		input    string
		wantOs   string
		wantArch string
		wantOk   bool
	}{
		{"gometalinter-2.0.5-darwin-amd64.tar.gz", "darwin", "amd64", true},
		{"gometalinter-2.0.5-windows-amd64.zip", "windows", "amd64", true},
		{"gometalinter-2.0.5-windows-amd64.zip", "windows", "amd64", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotOs, gotArch, gotOk := GetAssetMetadata(tt.input)
			if gotOs != tt.wantOs {
				t.Errorf("GetAssetMetadata() gotOs = %v, want %v", gotOs, tt.wantOs)
			}
			if gotArch != tt.wantArch {
				t.Errorf("GetAssetMetadata() gotArch = %v, want %v", gotArch, tt.wantArch)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetAssetMetadata() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestAssetDirectory(t *testing.T) {
	tests := []struct {
		parentDir, assetName, want string
	}{
		{"/foo/bar", "gometalinter-2.0.5-darwin-amd64.tar.gz", "/foo/bar/gometalinter-2.0.5-darwin-amd64"},
		{"/foo/bar", "gometalinter-2.0.5-windows-amd64.zip", "/foo/bar/gometalinter-2.0.5-windows-amd64"},
		//{"gometalinter-2.0.5-windows-amd64.zip", "windows", "amd64", true},
		//{"gometalinter-2.0.5-windows-amd64.zip", "windows", "amd64", true},
	}
	for _, tt := range tests {
		t.Run(tt.assetName, func(t *testing.T) {
			got := AssetDirectory(tt.parentDir, tt.assetName)
			if got != tt.want {
				t.Errorf("AssetDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZipFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"gometalinter-2.0.5-darwin-amd64.tar.gz", false},
		{"gometalinter-2.0.5-windows-amd64.zip", true},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := IsZipFile(tt.filename); got != tt.want {
				t.Errorf("IsZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

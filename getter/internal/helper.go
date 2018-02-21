package internal

import (
	"path/filepath"
	"regexp"
	"strings"
	"io"
)

// GetAssetMetadata gets some data from an asset filename
func GetAssetMetadata(name string) (os, arch string, ok bool) {
	assetMatcher := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-]+)(\.tar\.gz|\.tar\.bz2|\.zip)`)
	if !assetMatcher.MatchString(name) {
		return
	}
	m := assetMatcher.FindStringSubmatch(name)
	return m[3], m[4], true
}

// AssetDirectory builds the directory an asset should be extracted to
func AssetDirectory(parentDir, assetName string) string {
	re := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-^\.]+)(\.tar\.gz|\.tar\.bz2|\.zip)`)
	m := re.FindStringSubmatch(assetName)
	suffix := m[5]
	return filepath.Join(parentDir, strings.TrimSuffix(assetName, suffix))
}

// IsZipFile determines if a file ends with .zip
func IsZipFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".zip")
}

// SafeClose handles deferred errors
func SafeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

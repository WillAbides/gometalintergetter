package internal

import (
	"path/filepath"
	"regexp"
	"strings"
)

func GetAssetMetadata(name string) (os, arch string, ok bool) {
	assetMatcher := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-]+)(\.tar\.gz|\.tar\.bz2|\.zip)`)
	if !assetMatcher.MatchString(name) {
		return
	}
	m := assetMatcher.FindStringSubmatch(name)
	return m[3], m[4], true
}

func AssetDirectory(parentDir, assetName string) string {
	re := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-^\.]+)(\.tar\.gz|\.tar\.bz2|\.zip)`)
	m := re.FindStringSubmatch(assetName)
	suffix := m[5]
	return filepath.Join(parentDir, strings.TrimSuffix(assetName, suffix))
}

func IsZipFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".zip")
}

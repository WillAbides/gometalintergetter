package getter

import (
	"strings"
	"context"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"runtime"
	"regexp"
	"net/http"
	"github.com/mholt/archiver"
	"path/filepath"
	"io/ioutil"
)

var defaultRepo = repository{"alecthomas", "gometalinter"}

type (
	repository struct{ owner, name string }

	downloader struct {
		arch, os   string
		repoSvc    RepositorySvc
		repository *repository
		force      bool
	}

	Option func(*downloader)

	RepositorySvc interface {
		GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
		GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	}
)

// WithArch specifies the OS to download.  Default is the same as the current runtime
func WithOS(os string) Option {
	return func(d *downloader) {
		d.os = os
	}
}

// WithArch specifies the system architecture to download.  Default is the same as the current runtime
func WithArch(arch string) Option {
	return func(d *downloader) {
		d.arch = arch
	}
}

// WithMetalinterRepo specifies what GitHub repo gometalinter is in
func WithMetalinterRepo(repo string) Option {
	sl := strings.Split(repo, "/")
	r := &repository{}
	if len(sl) == 2 {
		r.owner, r.name = sl[0], sl[1]
	}
	return func(d *downloader) {
		d.repository = r
	}
}

// WithRepositoryService specifies the github.Repositories service to use
func WithRepositoryService(repoSvc RepositorySvc) Option {
	return func(d *downloader) {
		d.repoSvc = repoSvc
	}
}

// WithForce causes gometalinter to be downloaded even if it is already present
func WithForce() Option {
	return func(d *downloader) {
		d.force = true
	}
}

func getAssetMetadata(name string) (os, arch string, ok bool) {
	assetMatcher := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-]+).tar.[bg]z2?`)
	if ! assetMatcher.MatchString(name) {
		return
	}
	m := assetMatcher.FindStringSubmatch(name)
	return m[3], m[4], true
}

func (d *downloader) getRelease(version string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	var release *github.RepositoryRelease
	var err error
	if version == "" {
		release, _, err = d.repoSvc.GetLatestRelease(ctx, d.repository.owner, d.repository.name)
	} else {
		tag := "v" + version
		release, _, err = d.repoSvc.GetReleaseByTag(ctx, d.repository.owner, d.repository.name, tag)
	}
	return release, errors.Wrap(err, "failed to get release")
}

func (d *downloader) getReleaseAsset(version string) (*github.ReleaseAsset, error) {
	release, err := d.getRelease(version)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting release")
	}
	assets := release.Assets
	for _, asset := range assets {
		assetOS, assetArch, ok := getAssetMetadata(asset.GetName())
		if ok == false || assetArch != d.arch || assetOS != d.os {
			continue
		}
		return &asset, nil
	}
	return nil, errors.Errorf("could not find release asset")
}

func DownloadMetalinter(version, dstPath string, opts ...Option) error {
	d := &downloader{
		arch:       runtime.GOARCH,
		os:         runtime.GOOS,
		repository: &defaultRepo,
		repoSvc:    github.NewClient(nil).Repositories,
	}
	for _, opt := range opts {
		opt(d)
	}

	asset, err := d.getReleaseAsset(version)
	if err != nil {
		return errors.Wrap(err, "failed getting release asset")
	}

	archiveURLFile := filepath.Join(assetDirectory(dstPath, asset.GetName()), ".archiveurl")

	if ! d.force {
		oldUrl, err := ioutil.ReadFile(archiveURLFile)
		if err == nil && string(oldUrl) == asset.GetBrowserDownloadURL() {
			return nil
		}
	}

	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return errors.Wrapf(err, "failed downloading file: %s", asset.GetBrowserDownloadURL())
	}
	defer resp.Body.Close()

	archy := archiver.MatchingFormat(asset.GetName())

	err = archy.Read(resp.Body, dstPath)
	if err != nil {
		return errors.Wrap(err, "failed extracting archive")
	}

	err = ioutil.WriteFile(archiveURLFile, []byte(asset.GetBrowserDownloadURL()), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed writing %v", archiveURLFile)
	}
	return nil
}

func assetDirectory(parentDir, assetName string) string {
	re := regexp.MustCompile(`([^-]+)-v?(.+)-([^-]+)-([^-^\.]+)(.tar.[bg]z2?)?`)
	m := re.FindStringSubmatch(assetName)
	suffix := m[5]
	return filepath.Join(parentDir, strings.TrimSuffix(assetName, suffix))
}

package getter

import (
	"context"
	"github.com/WillAbides/gometalintergetter/getter/internal"
	"github.com/google/go-github/github"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var defaultRepo = repository{"alecthomas", "gometalinter"}

type (
	repository struct{ owner, name string }

	downloader struct {
		arch, os    string
		repoSvc     RepositorySvc
		repository  *repository
		force       bool
		skipSymlink bool
	}

	// Option is an configuration option for DownloadMetalinter
	Option func(*downloader)

	// RepositorySvc is the portion of github.Repositories that we use
	RepositorySvc interface {
		GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
		GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	}
)

// WithOS specifies the OS to download.  Default is the same as the current runtime
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

// SkipSymlink causes getter to not create a symlink for gometalinter
func SkipSymlink() Option {
	return func(d *downloader) {
		d.skipSymlink = true
	}
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
		assetOS, assetArch, ok := internal.GetAssetMetadata(asset.GetName())
		if !ok || assetArch != d.arch || assetOS != d.os {
			continue
		}
		return &asset, nil
	}
	return nil, errors.Errorf("could not find release asset")
}

// DownloadMetalinter download gometalinter to the specified path and make a symlink
func DownloadMetalinter(version, dstPath string, opts ...Option) error {
	var err error
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

	archiveURLFile := filepath.Join(internal.AssetDirectory(dstPath, asset.GetName()), ".archiveurl")
	binFile := filepath.Join(internal.AssetDirectory(dstPath, asset.GetName()), "gometalinter")

	if !d.force {
		oldURL, err := ioutil.ReadFile(archiveURLFile)
		if err == nil && string(oldURL) == asset.GetBrowserDownloadURL() {
			return nil
		}
	}

	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return errors.Wrapf(err, "failed downloading file: %s", asset.GetBrowserDownloadURL())
	}

	defer internal.SafeClose(resp.Body, &err)

	archy := archiver.MatchingFormat(asset.GetName())

	archiverTarget := dstPath
	if internal.IsZipFile(asset.GetName()) {
		archiverTarget = internal.AssetDirectory(dstPath, asset.GetName())
	}
	err = archy.Read(resp.Body, archiverTarget)
	if err != nil {
		return errors.Wrap(err, "failed extracting archive")
	}

	err = ioutil.WriteFile(archiveURLFile, []byte(asset.GetBrowserDownloadURL()), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed writing %v", archiveURLFile)
	}

	if !d.skipSymlink {
		err := os.Symlink(binFile, filepath.Join(dstPath, "gometalinter"))
		if err != nil {
			return errors.Wrap(err, "failed creating symlink")
		}
	}

	return err
}

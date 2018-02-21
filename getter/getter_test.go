package getter_test

import (
	"path/filepath"
	"testing"
	"net/http"
	. "github.com/WillAbides/gometalintergetter/getter"
	"github.com/WillAbides/gometalintergetter/getter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/google/go-github/github"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
)

type (
	testData struct {
		*testing.T
		workDir string
	}
)

func testSetup(t *testing.T) *testData {
	t.Helper()
	wd, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	return &testData{
		workDir: wd,
		T:       t,
	}
}

func (td *testData) teardown() {
	td.Helper()
	err := os.RemoveAll(td.workDir)
	if err != nil {
		td.Fatal(err)
	}
}

func TestDownloadMetalinter(t *testing.T) {
	td := testSetup(t)
	defer td.teardown()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/gometalinter-v2.0.2-linux-amd64.tar.gz")
	}))
	defer ts.Close()

	aName := "gometalinter-v2.0.2-linux-amd64.tar.gz"
	aDownloadUrl := ts.URL
	mockRelease := &github.RepositoryRelease{
		Assets: []github.ReleaseAsset{{
			Name:               &aName,
			BrowserDownloadURL: &aDownloadUrl,
		}},
	}

	repoSvc := new(mocks.RepositorySvc)
	repoSvc.On("GetReleaseByTag", mock.Anything, "alecthomas", "gometalinter", "v2.0.2").Return(mockRelease, nil, nil)

	err := DownloadMetalinter("2.0.2", td.workDir,
		WithRepositoryService(repoSvc),
		WithOS("linux"),
		WithArch("amd64"),
	)
	assert.Nil(t, err)
	ar, err := ioutil.ReadFile(filepath.Join(td.workDir, "gometalinter-v2.0.2-linux-amd64",  ".archiveurl"))
	assert.Nil(t, err)
	assert.Equal(t, ts.URL, string(ar))
}


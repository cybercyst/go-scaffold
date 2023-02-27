package download

import (
	"testing"

	"github.com/spf13/afero"
	"gotest.tools/v3/assert"
)

func testGitUrl(url string, t *testing.T) {
	got := isGitRepo(url)
	if got != true {
		t.Errorf("%s was not detected as a valid git repository", url)
	}
}

func TestDetectIsGitRepositoryShouldDetectValidRepositories(t *testing.T) {
	// support https
	testGitUrl("https://github.com/user/template.git", t)

	// support SSH
	testGitUrl("git@github.com:user/template.git", t)

	// support self hosted repos
	testGitUrl("git@myurltest.com:user/template.git", t)
}

func TestDetectIsOciArtifactUriShouldDetectValidUri(t *testing.T) {
	got := isOciArtifactUri("oci://registry.url/repo/artifact:tag")
	assert.Equal(t, got, true, "oci://registry.url/repo/artifact:tag was not detected as a valid OCI artifact uri")
}

func TestDectectIsDirectoryShouldDetectValidDirectory(t *testing.T) {
	got := isDirectory(".")
	assert.Equal(t, got, true, "current test directory was not detected as a valid directory")

	got = isDirectory("i-dont-exist")
	assert.Equal(t, got, false, "non existant directory was detected as a valid directory")
}

func TestDetectErrorWhenNoValidUriPassed(t *testing.T) {
	fs := afero.NewMemMapFs()
	_, got := Download(fs, "this-isn't-a-valid-uri-or-folder")
	assert.Error(t, got, "open this-isn't-a-valid-uri-or-folder: file does not exist")
}

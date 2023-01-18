package download

import (
	"testing"
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
	if got != true {
		t.Error("oci://registry.url/repo/artifact:tag was not detected as a valid OCI artifact uri")
	}
}

func TestDectectIsDirectoryShouldDetectValidDirectory(t *testing.T) {
	got := isDirectory(".")
	if got != true {
		t.Error("current test directory was not detected as a valid directory")
	}

	got = isDirectory("i-dont-exist")
	if got == true {
		t.Error("non existant directory was detected as a valid directory")
	}
}

func TestDetectErrorWhenNoValidUriPassed(t *testing.T) {
	_, got := Download("this-isn't-a-valid-uri-or-folder")
	if got == nil {
		t.Error("invalid uri did not cause an error")
	}
}

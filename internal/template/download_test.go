package template

import "testing"

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

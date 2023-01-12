package internal

import "testing"

func TestDetectIsGitRepositoryShouldDetectValidRepositories(t *testing.T) {
	got := isGitRepo("https://github.com/user/repo.git")
	if got != true {
		t.Error("https://github.com/user/repo.git was not detected as a valid git repository")
	}

	got = isGitRepo("git@github.com:cybercyst/go-api.git")
	if got != true {
		t.Error("git@github.com:cybercyst/go-api.git was not detected as a valid git repository")
	}
}

func TestDetectIsOrasArtifactUriShouldDetectValidUri(t *testing.T) {
	got := isOrasArtifactUri("oci://registry.url/repo/artifact:tag")
	if got != true {
		t.Error("oci://registry.url/repo/artifact:tag was not detected as a valid OCI artifact uri")
	}
}

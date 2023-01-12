package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

func Download(uri string) (string, error) {
	switch {
	case isOrasArtifactUri(uri):
		return downloadOras(uri)
	case isGitRepo(uri):
		return downloadGit(uri)
	default:
		return downloadFs(uri)
	}
}

func createTempDir() string {
	path, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s-", ProgramName))
	if err != nil {
		log.Fatal("Unable to create temporary directory: ", err)
	}

	return path
}

func isOrasArtifactUri(uri string) bool {
	isOrasRegExp := regexp.MustCompile("^oci://")
	return isOrasRegExp.MatchString(uri)
}

func isGitRepo(uri string) bool {
	isSchemeRegExp := regexp.MustCompile(`^[^:]+://`)
	return isSchemeRegExp.MatchString(uri)
}

func downloadGit(gitRepo string) (string, error) {
	tempDir := createTempDir()

	_, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: gitRepo,
	})
	return tempDir, err
}

func downloadOras(uri string) (string, error) {
	ctx := context.Background()

	repoUri := strings.ReplaceAll(uri, "oci://", "")
	repo, err := remote.NewRepository(repoUri)
	if err != nil {
		return "", err
	}

	tempDir := createTempDir()
	dst := file.New(tempDir)

	copyOptions := oras.DefaultCopyOptions
	desc, err := oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, copyOptions)
	if err != nil {
		return "", err
	}

	fmt.Println("tempDir", tempDir)
	fmt.Println("Digest:", desc.Digest)
	return tempDir, nil
}

func downloadFs(path string) (string, error) {
	_, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

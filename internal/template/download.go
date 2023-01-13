package template

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cybercyst/go-cookiecutter/internal"
	"github.com/go-git/go-git/v5"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

func (t *Template) Download(uri string) error {
	t.Uri = uri

	switch {
	case isOciArtifactUri(t.Uri):
		return t.downloadOci()
	case isGitRepo(t.Uri):
		return t.downloadGit()
	default:
		return t.downloadFs()
	}
}

func createTempDir() string {
	path, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s-", internal.ProgramName))
	if err != nil {
		log.Fatal("Unable to create temporary directory: ", err)
	}

	return path
}

func isOciArtifactUri(uri string) bool {
	isOrasRegExp := regexp.MustCompile("^oci://")
	return isOrasRegExp.MatchString(uri)
}

func isGitRepo(uri string) bool {
	isGitRegExp := regexp.MustCompile(`^((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)
	return isGitRegExp.MatchString(uri)
}

func (t *Template) downloadGit() error {
	tempDir := createTempDir()

	repo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: t.Uri,
	})
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	t.Version = head.Hash().String()
	t.LocalPath = tempDir

	return err
}

func (t *Template) downloadOci() error {
	ctx := context.Background()

	repoUri := strings.ReplaceAll(t.Uri, "oci://", "")
	repo, err := remote.NewRepository(repoUri)
	if err != nil {
		return err
	}

	tempDir := createTempDir()
	dst := file.New(tempDir)

	copyOptions := oras.DefaultCopyOptions
	desc, err := oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, copyOptions)
	if err != nil {
		return err
	}

	t.Version = desc.Digest.String()
	t.LocalPath = tempDir

	return nil
}

func (t *Template) downloadFs() error {
	_, err := os.ReadDir(t.Uri)
	if err != nil {
		return err
	}

	t.Version = "HEAD"
	t.LocalPath = t.Uri

	return nil
}

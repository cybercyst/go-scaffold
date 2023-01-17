package template

import (
	"context"
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
		return t.useDirectory()
	}
}

func isOciArtifactUri(uri string) bool {
	isOrasRegExp := regexp.MustCompile("^oci://")
	return isOrasRegExp.MatchString(uri)
}

func (t *Template) downloadOci() error {
	ctx := context.Background()

	repoUri := strings.ReplaceAll(t.Uri, "oci://", "")
	repo, err := remote.NewRepository(repoUri)
	if err != nil {
		return err
	}

	tempDir := internal.CreateTempDir()
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

func isGitRepo(uri string) bool {
	isGitRegExp := regexp.MustCompile(`^((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)
	return isGitRegExp.MatchString(uri)
}

func (t *Template) downloadGit() error {
	tempDir := internal.CreateTempDir()

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

func isDirectory(uri string) bool {
	info, err := os.Stat(uri)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func (t *Template) useDirectory() error {
	_, err := os.ReadDir(t.Uri)
	if err != nil {
		return err
	}

	t.Version = "HEAD"
	t.LocalPath = t.Uri

	return nil
}

package download

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/cybercyst/go-cookiecutter/internal/utils"
)

type DownloadInfo struct {
	LocalPath string
	Version   string
}

func Download(uri string) (*DownloadInfo, error) {
	switch {
	case isOciArtifactUri(uri):
		return downloadOci(uri)
	case isGitRepo(uri):
		return downloadGit(uri)
	default:
		return useDirectory(uri)
	}
}

func isOciArtifactUri(uri string) bool {
	isOrasRegExp := regexp.MustCompile("^oci://")
	return isOrasRegExp.MatchString(uri)
}

func downloadOci(ociArtifactUri string) (*DownloadInfo, error) {
	ctx := context.Background()

	repoUri := strings.ReplaceAll(ociArtifactUri, "oci://", "")
	repo, err := remote.NewRepository(repoUri)
	if err != nil {
		return nil, err
	}

	tempDir := utils.CreateTempDir()
	dst := file.New(tempDir)

	copyOptions := oras.DefaultCopyOptions
	desc, err := oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, copyOptions)
	if err != nil {
		return nil, err
	}

	return &DownloadInfo{
		LocalPath: tempDir,
		Version:   desc.Digest.String(),
	}, nil
}

func isGitRepo(uri string) bool {
	isGitRegExp := regexp.MustCompile(`^((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)
	return isGitRegExp.MatchString(uri)
}

func downloadGit(gitRepo string) (*DownloadInfo, error) {
	tempDir := utils.CreateTempDir()

	repo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: gitRepo,
	})
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &DownloadInfo{
		LocalPath: tempDir,
		Version:   head.Hash().String(),
	}, nil
}

func isDirectory(uri string) bool {
	info, err := os.Stat(uri)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func useDirectory(filePath string) (*DownloadInfo, error) {
	if _, err := os.ReadDir(filePath); err != nil {
		return nil, err
	}

	// TODO: check if filePath is just a local git repo and
	// use its SHA for version

	return &DownloadInfo{
		LocalPath: filePath,
		Version:   "",
	}, nil
}

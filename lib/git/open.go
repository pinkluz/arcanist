package git

import (
	gogit "github.com/go-git/go-git/v5"
)

func OpenRepo() (*gogit.Repository, error) {
	repo, err := gogit.PlainOpenWithOptions(".", &gogit.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})

	return repo, err
}

package repo

import "github.com/mesuutt/git-mirror/pkg/commit"

type Repo interface {
	AddStats(diff *commit.Diff) error
	AddAndCommit(msg string) error
}

package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ValidateRepo(path string) error {
	if f, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s is not a git repository", path)
		}
		return err

	} else if !f.IsDir() {
		return fmt.Errorf("given stats repo path should be a directory. `%s` is not a directory", path)
	}

	return nil
}

func AddAndCommit(repoPath, msg string) error {
	_, err := exec.Command("git", "-C", repoPath, "add", ".").Output()
	if err != nil {
		return fmt.Errorf("`git -C %s add .` command execution failed. error %v", repoPath, err)
	}

	_, err = exec.Command("git", "-C", repoPath, "commit", "-m", msg).Output()
	if err != nil {
		return fmt.Errorf("`git -C %s commit -m '%s'` command execution failed. error %v", repoPath, msg, err)
	}

	return nil
}

func LastCommitStatsDiff(repoPath string) ([]byte, error) {
	out, err := exec.Command("git", "-C", repoPath, "diff", "--numstat", "HEAD~1").Output()
	if err != nil {
		return nil, fmt.Errorf("`git -C %s diff --numstat HEAD~1` command execution failed with error: `%v`", repoPath, err)
	}

	return out, nil
}

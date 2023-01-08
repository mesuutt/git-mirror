package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Empty tree hash
// https://github.com/git/git/blob/d9a3764af74ac215e06543c263ec21196d672b49/cache.h#L1027-L1028
const EmptyTreeHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

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

func AddChanges(repoPath string) error {
	_, err := exec.Command("git", "-C", repoPath, "add", ".").Output()
	if err != nil {
		return fmt.Errorf("`git -C %s add .` command execution failed. error %v", repoPath, err)
	}
	return nil
}

func Commit(path, msg string) error {
	_, err := exec.Command("git", "-C", path, "commit", "-m", msg).Output()
	if err != nil {
		return fmt.Errorf("`git -C %s commit -m '%s'` command execution failed. error %v", path, msg, err)
	}

	return nil
}

func LastCommitStatsDiff(path string) ([]byte, error) {
	revision, err := findDiffRevision(path)
	if err != nil {
		return nil, err
	}

	return numStatDiff(path, revision)
}

func numStatDiff(path, revision string) ([]byte, error) {
	out, err := exec.Command("git", "-C", path, "diff", "--numstat", revision).Output()
	if err != nil {
		return nil, fmt.Errorf("`git -C %s diff --numstat %s` command execution failed with error: `%v`", path, revision, err)
	}

	return out, nil
}

func findDiffRevision(path string) (string, error) {
	c, err := commitCount(path)
	if err != nil {
		return "", err
	}

	if c > 1 {
		return "HEAD~1", nil
	} else if c == 1 {
		return EmptyTreeHash, nil
	}

	return "", fmt.Errorf("git repo needs to have at least one commit")
}

// commitCount returns num 0, 1 or 2
func commitCount(path string) (int, error) {
	var out []byte
	var err error
	rev := "HEAD"

	// We don't need to count all commits, so we are giving --max-count=2
	out, err = exec.Command("git", "-C", path, "rev-list", "--count", "--max-count=2", rev).Output()
	if err != nil {
		// commit count less than 2
		// so checking there is any commit
		out, err = exec.Command("git", "-C", path, "rev-list", "--count", "--max-count=2", EmptyTreeHash).Output()
		if err != nil {
			return 0, errors.New("git commit count getting failed")
		}
	}

	c, err := strconv.Atoi(strings.Replace(string(out), "\n", "", -1))
	if err != nil {
		return 0, fmt.Errorf("commit count parse failed")
	}

	return c, nil
}

// AddHook adds given content to the given hook file
func AddHook(path string, hook, content string) error {
	if err := createHookDir(path); err != nil {
		return err
	}

	// TODO: prevent duplicate hook
	// TODO: in test check file is executable
	hookFilePath := filepath.Join(path, ".git", "hooks", hook)
	f, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%s file could not open for adding hook: %v", hookFilePath, err)
	}

	defer f.Close()

	if _, err := f.WriteString("\n" + content + "\n"); err != nil {
		return fmt.Errorf("hook add failed to %s file. error: %v", hook, err)
	}

	return nil
}

func createHookDir(path string) error {
	hookDir := filepath.Join(path, ".git", "hooks")
	if _, err := os.Stat(hookDir); os.IsNotExist(err) {
		if err := os.MkdirAll(hookDir, os.ModeDir); err != nil {
			return fmt.Errorf("%s directory creation failed: %v", hookDir, err)
		}
	}

	return nil
}

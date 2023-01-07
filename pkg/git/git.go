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

func AddChanges(repoPath string) error {
	_, err := exec.Command("git", "-C", repoPath, "add", ".").Output()
	if err != nil {
		return fmt.Errorf("`git -C %s add .` command execution failed. error %v", repoPath, err)
	}
	return nil
}

func Commit(repoPath, msg string) error {
	_, err := exec.Command("git", "-C", repoPath, "commit", "-m", msg).Output()
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

func CreateHookDir(path string) error {
	hookDir := filepath.Join(path, ".git", "hooks")
	if _, err := os.Stat(hookDir); os.IsNotExist(err) {
		if err := os.MkdirAll(hookDir, os.ModeDir); err != nil {
			return fmt.Errorf(hookDir+" directory creation failed: %v", err)
		}
	}

	return nil
}

func AddHook(path string, hook, content string) error {
	// TODO: prevent duplicate hook
	// TODO: in test check file is executable
	hookFilePath := filepath.Join(path, ".git", "hooks", hook)
	f, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf(hookFilePath+" file could not open for adding hook: %v", err)
	}

	defer f.Close()

	if _, err := f.WriteString("\n" + content + "\n"); err != nil {
		return fmt.Errorf("hook add failed to post-commit file. error: %v", err)
	}

	return nil
}

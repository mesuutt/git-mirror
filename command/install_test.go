package command

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestInstallHookCmd(t *testing.T) {
	fs := afero.NewMemMapFs()
	ctx := cli.NewContext(&cli.App{Writer: io.Discard}, nil, nil)

	ctx.Context = context.WithValue(ctx.Context, "fs", fs)

	hookDir := filepath.Join(".git", "hooks")
	hookFile := filepath.Join(hookDir, "post-commit")
	_ = fs.MkdirAll(hookDir, os.ModeDir)

	// when
	command := cli.Command{Action: InstallHookCmd}
	err := command.Run(ctx, []string{"install"}...)

	// then
	assert.Nil(t, err)
	contains, err := afero.FileContainsBytes(fs, hookFile, []byte("git-mirror add"))

	assert.Nil(t, err)
	assert.True(t, contains)
}

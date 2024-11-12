package commands_test

import (
	"os"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	"github.com/stretchr/testify/require"
)

func TestClone(t *testing.T) {
	t.Run("clone correctly some repo", func(t *testing.T) {
		args := []string{"https://github.com/Luisgustavom1/Luisgustavom1.git", "my-repo"}

		res, err := commands.Clone(args)
		require.NoError(t, err)
		require.Equal(t, "", res)

		_, err = os.Stat("my-repo")
		require.False(t, os.IsNotExist(err), "dir my-repo should be created")
	})

	t.Run("invalid arguments", func(t *testing.T) {
		args := []string{"https://some-url"}

		_, err := commands.Clone(args)
		require.EqualErrorf(t, err, "usage: mygit clone <url> <some_dir>\n", "Invalid args")
	})
}

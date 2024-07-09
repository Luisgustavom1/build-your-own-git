package commands_test

import (
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	"github.com/stretchr/testify/require"
)

func TestCommitTree(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "commit tree with no flags",
			args:     []string{},
			expected: "No tree hash provided",
		},
		{
			name:     "commit tree without parent hash",
			args:     []string{"tree_hash", "-p"},
			expected: "No parent hash provided",
		},
		{
			name:     "commit tree without message",
			args:     []string{"tree_hash", "-m"},
			expected: "No message provided",
		},
		{
			name:     "commit tree without message",
			args:     []string{"tree_hash", "-p", "parent_hash", "-m"},
			expected: "No message provided",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := commands.CommitTree(tc.args)
			require.NotNil(t, err)
			require.Equal(t, tc.expected, err.Error())
			require.Equal(t, "", res)
		})
	}
}

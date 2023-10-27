package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envDirPath := "./testdata/env"
	tests := []struct {
		name     string
		cmd      []string
		exitCode int
	}{
		{
			name: "test showing files on the dir",
			cmd:  []string{"ls", "-la", "/home"},
		},
		{
			name: "test showing file size",
			cmd:  []string{"du", "-h", "/usr/lib/os-release"},
		},
		{
			name:     "test trying showing protected files",
			cmd:      []string{"ls", "-la", "/root"},
			exitCode: 2,
		},
	}

	environments, err := ReadDir(envDirPath)
	if err != nil {
		require.Fail(t, "Error: %v", err)
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exitCode := RunCmd(tc.cmd, environments)
			require.Equal(t, tc.exitCode, exitCode)
		})
	}
}

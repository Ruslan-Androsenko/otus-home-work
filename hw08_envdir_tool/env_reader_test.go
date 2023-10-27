package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envDirPath := "./testdata/env"
	existsVariables := []string{"BAR", "EMPTY", "FOO", "HELLO"}
	notExistsVariables := []string{"UNSET"}

	t.Run("check exists environment variables", func(t *testing.T) {
		environments, err := ReadDir(envDirPath)
		if err != nil {
			require.Fail(t, "Error: %v", err)
		}

		require.NoError(t, err)
		require.NotNil(t, environments)

		for _, key := range existsVariables {
			valueMap, okMap := environments[key]

			require.True(t, okMap)
			require.False(t, valueMap.NeedRemove)
		}
	})

	t.Run("check not exists environment variables", func(t *testing.T) {
		environments, err := ReadDir(envDirPath)
		if err != nil {
			require.Fail(t, "Error: %v", err)
		}

		require.NoError(t, err)
		require.NotNil(t, environments)

		for _, key := range notExistsVariables {
			valueMap, okMap := environments[key]

			require.True(t, okMap)
			require.True(t, valueMap.NeedRemove)
			require.Empty(t, valueMap.Value)
		}
	})
}

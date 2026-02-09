package file_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anchore/go-make/file"
	"github.com/anchore/go-make/require"
)

func Test_Ls(t *testing.T) {
	ls := file.Ls("testdata/some/.other")
	require.Contains(t, ls, ".config.json")
	require.Contains(t, ls, ".config.yaml")

	// Test with temp directory with known files
	tmp := t.TempDir()
	file1 := filepath.Join(tmp, "file1.txt")
	file2 := filepath.Join(tmp, "file2.txt")
	require.NoError(t, os.WriteFile(file1, []byte("content1"), 0o644))
	require.NoError(t, os.WriteFile(file2, []byte("content2"), 0o755))

	ls = file.Ls(tmp)
	require.Contains(t, ls, "file1.txt")
	require.Contains(t, ls, "file2.txt")
	require.Contains(t, ls, "-rw-r--r--") // permissions for file1
	require.Contains(t, ls, "-rwxr-xr-x") // permissions for file2
}

func Test_LsAll(t *testing.T) {
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "subdir")
	require.NoError(t, os.Mkdir(subdir, 0o755))

	file1 := filepath.Join(tmp, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	require.NoError(t, os.WriteFile(file1, []byte("content1"), 0o644))
	require.NoError(t, os.WriteFile(file2, []byte("content2"), 0o644))

	lsAll := file.LsAll(tmp)

	// Should contain root level file
	require.Contains(t, lsAll, "file1.txt")
	// Should contain subdirectory indicator
	require.Contains(t, lsAll, "subdir:")
	// Should contain file in subdirectory
	require.Contains(t, lsAll, "file2.txt")
}

func Test_LsAll_Testdata(t *testing.T) {
	lsAll := file.LsAll("testdata/some/.other")
	require.Contains(t, lsAll, ".config.json")
	require.Contains(t, lsAll, ".config.yaml")
}

func Test_LogWorkdir(t *testing.T) {
	tmp := t.TempDir()
	file.InDir(tmp, func() {
		file.Write(filepath.Join(tmp, "file.txt"), "hello world")
		file.LogWorkdir()
	})
}

func Test_HumanizeBytes(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "bytes",
			size:     512,
			expected: "512",
		},
		{
			name:     "kilobytes",
			size:     2048,
			expected: "2KB",
		},
		{
			name:     "megabytes",
			size:     5 * 1024 * 1024,
			expected: "5MB",
		},
		{
			name:     "gigabytes",
			size:     3 * 1024 * 1024 * 1024,
			expected: "3GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := file.HumanizeBytes(tt.size)
			resultStr, ok := result.(string)
			require.True(t, ok)
			require.True(t, strings.Contains(resultStr, tt.expected))
		})
	}
}

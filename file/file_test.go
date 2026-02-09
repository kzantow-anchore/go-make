package file_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anchore/go-make/file"
	"github.com/anchore/go-make/lang"
	"github.com/anchore/go-make/require"
)

func Test_Cwd(t *testing.T) {
	tests := []struct {
		name string
		dir  string
	}{
		{
			name: "current directory",
			dir:  ".",
		},
		{
			name: "other directory",
			dir:  "testdata",
		},
	}

	startDir := lang.Return(os.Getwd())
	defer func() { require.NoError(t, os.Chdir(startDir)) }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file.InDir(tt.dir, func() {
				expected := lang.Return(filepath.Abs(filepath.Join(startDir, tt.dir)))
				got := lang.Return(filepath.Abs(file.Cwd()))
				require.Equal(t, expected, got)
			})
		})
	}
}

func Test_Copy(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "srcfile.txt")
	require.NoError(t, os.WriteFile(srcFile, []byte("hello world"), 0o742))
	srcPerms := lang.Return(os.Stat(srcFile)).Mode()
	dstFile := filepath.Join(tempDir, "tmpdir", "dstfile.txt")

	file.Copy(srcFile, dstFile)

	require.Equal(t, "hello world", file.Read(dstFile))
	perms := lang.Return(os.Stat(dstFile)).Mode()
	require.Equal(t, srcPerms, perms.Perm())
}

func Test_IsRegular(t *testing.T) {
	tempDir := t.TempDir()
	require.True(t, !file.IsRegular(tempDir))

	srcFile := filepath.Join(tempDir, "srcfile.txt")
	require.NoError(t, os.WriteFile(srcFile, []byte("hello world"), 0o742))
	require.True(t, file.IsRegular(srcFile))
}

func Test_FindParent(t *testing.T) {
	tests := []struct {
		file     string
		expected string
	}{
		{
			file:     ".config.yaml",
			expected: "some/.config.yaml",
		},
		{
			file:     ".config.json",
			expected: "some/nested/path/.config.json",
		},
		{
			file:     ".other",
			expected: "some/.other",
		},
		{
			file:     ".missing",
			expected: "",
		},
	}
	testdataDir, _ := os.Getwd()
	testdataDir = filepath.ToSlash(filepath.Join(testdataDir, "testdata")) + "/"
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			file.InDir("testdata/some/nested/path", func() {
				path := file.FindParent(file.Cwd(), tt.file)
				path = filepath.ToSlash(path)
				path = strings.TrimPrefix(path, testdataDir)
				require.Equal(t, tt.expected, path)
			})
		})
	}
}

func Test_EnsureDir(t *testing.T) {
	testDir := t.TempDir()

	newDir := filepath.Join(testDir, "newdir")

	require.True(t, !file.Exists(newDir))

	file.EnsureDir(newDir)

	require.True(t, file.Exists(newDir))

	// existing dir, does nothing, no error
	file.EnsureDir(newDir)

	require.True(t, file.Exists(newDir))

	newFile := filepath.Join(testDir, "newfile")
	require.NoError(t, os.WriteFile(newFile, []byte("hello world"), 0o700))

	// should panic if unable to create dir such as when a file exists
	require.Error(t, lang.Catch(func() {
		file.EnsureDir(newFile)
	}))
}

func Test_FindAll(t *testing.T) {
	tests := []struct {
		pattern       string
		expectedCount int
	}{
		{
			pattern:       "**/*.json",
			expectedCount: 2,
		},
		{
			pattern:       "**/*",
			expectedCount: 4,
		},
		{
			pattern:       ".config.yaml",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			file.InDir("testdata/some", func() {
				got := file.FindAll(tt.pattern)
				require.Equal(t, tt.expectedCount, len(got))
			})
		})
	}
}

func Test_ChmodAll_SingleFile(t *testing.T) {
	tmp := t.TempDir()
	testFile := filepath.Join(tmp, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0o600))

	err := file.ChmodAll(testFile, 0o755)
	require.NoError(t, err)

	stat := lang.Return(os.Stat(testFile))
	require.Equal(t, os.FileMode(0o755), stat.Mode().Perm())
}

func Test_ChmodAll_DirectoryWithNestedFiles(t *testing.T) {
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "subdir")
	require.NoError(t, os.Mkdir(subdir, 0o700))

	file1 := filepath.Join(tmp, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	require.NoError(t, os.WriteFile(file1, []byte("content1"), 0o600))
	require.NoError(t, os.WriteFile(file2, []byte("content2"), 0o600))

	err := file.ChmodAll(tmp, 0o755)
	require.NoError(t, err)

	// Check directory permissions
	stat := lang.Return(os.Stat(tmp))
	require.Equal(t, os.FileMode(0o755)|os.ModeDir, stat.Mode())

	// Check subdirectory permissions
	stat = lang.Return(os.Stat(subdir))
	require.Equal(t, os.FileMode(0o755)|os.ModeDir, stat.Mode())

	// Check file permissions
	stat = lang.Return(os.Stat(file1))
	require.Equal(t, os.FileMode(0o755), stat.Mode().Perm())

	stat = lang.Return(os.Stat(file2))
	require.Equal(t, os.FileMode(0o755), stat.Mode().Perm())
}

func Test_ChmodAll_DeeplyNestedStructure(t *testing.T) {
	tmp := t.TempDir()
	dir1 := filepath.Join(tmp, "dir1")
	dir2 := filepath.Join(dir1, "dir2")
	dir3 := filepath.Join(dir2, "dir3")
	require.NoError(t, os.MkdirAll(dir3, 0o700))

	deepFile := filepath.Join(dir3, "deep.txt")
	require.NoError(t, os.WriteFile(deepFile, []byte("deep content"), 0o600))

	err := file.ChmodAll(tmp, 0o755)
	require.NoError(t, err)

	// Check all directories have correct permissions
	for _, dir := range []string{tmp, dir1, dir2, dir3} {
		stat := lang.Return(os.Stat(dir))
		require.Equal(t, os.FileMode(0o755)|os.ModeDir, stat.Mode())
	}

	// Check file has correct permissions
	stat := lang.Return(os.Stat(deepFile))
	require.Equal(t, os.FileMode(0o755), stat.Mode().Perm())
}

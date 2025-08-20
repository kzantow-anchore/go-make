package github

import (
	"path/filepath"
	"testing"

	"github.com/anchore/go-make/config"
	"github.com/anchore/go-make/file"
	"github.com/anchore/go-make/log"
	"github.com/anchore/go-make/require"
	"github.com/anchore/go-make/run"
	"github.com/anchore/go-make/script"
)

func Test_UploadDownload(t *testing.T) {
	if !config.CI {
		t.Log("skipping artifact upload test in non-CI environment")
		return
	}

	tests := []struct {
		files    string
		artifact string
	}{
		{
			files:    "pr_run_artifacts.json",
			artifact: "my-artifact-name",
		},
		{
			files: "**/*_artifacts.json",
		},
		{
			files: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.files, func(t *testing.T) {
			p := Payload() // tests run in workflow in github

			api := NewClient(Owner("testorg"), Repo("testrepo"))
			artifactId := api.UploadArtifactDir("testdata", UploadArtifactOption{
				ArtifactName: tt.artifact,
				Glob:         tt.files,
			})
			require.True(t, artifactId != 0)

			tmpdir := t.TempDir()

			api.DownloadArtifactDir(HeadSha(p.SHA), p.Workflow, tt.artifact, tmpdir)

			testdataFile := file.Read("testdata/pr_run_artifacts.json")
			downloadedFile := file.Read(filepath.Join(tmpdir, "pr_run_artifacts.json"))

			require.Equal(t, testdataFile, downloadedFile)

			require.True(t, file.Exists(filepath.Join(tmpdir, "testdata/empty.json")))
			require.True(t, !file.Exists(filepath.Join(tmpdir, "empty.json")))
		})
	}
}

func Test_ensureActionsArtifactNpmPackageInstalled(t *testing.T) {
	if !config.CI {
		t.Log("skipping artifact upload test in non-CI environment")
		return
	}

	ensureActionsArtifactInstalled()

	log.Log("actions/artifact out: %s", script.Run("npm list -g @actions/artifact", run.NoFail()))
}

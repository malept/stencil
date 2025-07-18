package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaredallard/cmdexec"
	"github.com/urfave/cli/v3"
	"go.rgst.io/stencil/v2/pkg/configuration"
	"go.rgst.io/stencil/v2/pkg/slogext"
	"gotest.tools/v3/assert"
)

// prepareTestRun sets up the environment for running a stencil command.
//
// dir is the directory to change into before running the command. If
// dir is an empty string, a temporary directory will be created.
func prepareTestRun(t *testing.T, dir string) {
	// Change into the repo root.
	// GOMOD: <module path>/go.mod
	b, err := cmdexec.Command("go", "env", "GOMOD").Output()
	assert.NilError(t, err)
	repoRoot := strings.TrimSuffix(strings.TrimSpace(string(b)), "/go.mod")

	t.Chdir(repoRoot)

	// Use a temporary directory for the test if one is not provided.
	if dir == "" {
		dir = t.TempDir()
	}
	t.Chdir(dir)
}

// testRunCommand runs the provided [cli.Command].
func testRunCommand(t *testing.T, cmd *cli.Command, dir string, args ...string) error {
	prepareTestRun(t, dir)

	wCmd := &cli.Command{Name: "test", Commands: []*cli.Command{cmd}}
	return wCmd.Run(t.Context(), append([]string{"test", cmd.Name}, args...))
}

func TestCanCreateModule(t *testing.T) {
	log := slogext.NewTestLogger(t)
	cmd := NewModuleCreateCommand(log)
	assert.Assert(t, cmd != nil)
	assert.NilError(t, testRunCommand(t, cmd, "", "github.com/rgst-io/test-module"))

	// Ensure it created the expected files.
	_, err := os.Stat("stencil.yaml")
	assert.NilError(t, err)
}

func TestCreateModuleFailsWhenFilesExist(t *testing.T) {
	log := slogext.NewTestLogger(t)
	cmd := NewModuleCreateCommand(log)
	assert.Assert(t, cmd != nil)

	tmpDir := t.TempDir()

	// Create a file to trigger the error.
	f, err := os.Create(filepath.Join(tmpDir, "test-file"))
	assert.NilError(t, err)
	assert.NilError(t, f.Close())

	err = testRunCommand(t, cmd, tmpDir, "github.com/rgst-io/test-module")
	assert.ErrorContains(t, err, "directory is not empty, found test-file")
}

// TestCanCreateNativeExtension ensures that we can render a native
// extension through stencil-golang. This is technically more of a test
// of stencil-golang than anything else.
func TestCanCreateNativeExtension(t *testing.T) {
	log := slogext.NewTestLogger(t)
	cmd := NewModuleCreateCommand(log)
	assert.Assert(t, cmd != nil)

	assert.NilError(t, testRunCommand(t, cmd, "", "--native-extension", "github.com/rgst-io/test-module"))

	// Ensure it created the expected files.
	expectedFiles := []string{
		filepath.Join("cmd", "plugin", "plugin.go"),
		"stencil.yaml",
	}
	for _, f := range expectedFiles {
		_, err := os.Stat(f)
		assert.NilError(t, err)
	}

	// Ensure that we have 'plugin' as one of our types.
	tr, err := configuration.LoadDefaultTemplateRepositoryManifest()
	assert.NilError(t, err)
	assert.Assert(t, tr.Type.Contains(configuration.TemplateRepositoryTypeExt))
}

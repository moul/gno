package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExamplesTxtar(t *testing.T) {
	t.Parallel()

	// Find all .txtar files in the examples directory tree
	examplesDir := filepath.Join("..", "..", "..", "examples")
	var txtarFiles []string
	
	err := filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".txtar") {
			txtarFiles = append(txtarFiles, path)
		}
		return nil
	})
	require.NoError(t, err)

	if len(txtarFiles) == 0 {
		t.Skip("no .txtar files found in examples directory")
	}

	// Run each .txtar file as a separate test
	for _, txtarFile := range txtarFiles {
		txtarFile := txtarFile // capture for closure
		relPath, err := filepath.Rel(examplesDir, txtarFile)
		require.NoError(t, err)
		
		testName := strings.TrimSuffix(relPath, ".txtar")
		testName = strings.ReplaceAll(testName, string(filepath.Separator), "_")
		
		t.Run(testName, func(t *testing.T) {
			// Create a temporary directory and copy the .txtar file there
			tmpDir := t.TempDir()
			dstPath := filepath.Join(tmpDir, filepath.Base(txtarFile))
			
			// Read and write the file
			content, err := os.ReadFile(txtarFile)
			require.NoError(t, err)
			err = os.WriteFile(dstPath, content, 0644)
			require.NoError(t, err)
			
			// Run the testscripts using the helper
			RunGnolandTestscripts(t, tmpDir)
		})
	}
}
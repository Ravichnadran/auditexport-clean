package output

import (
	"archive/zip"
	"auditexport/internal/run"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipEvidence() error {
	// Create zip at run root
	zipPath := filepath.Join(run.BaseDir(), "evidence.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	evidenceRoot := run.EvidencePath()

	return filepath.Walk(evidenceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and the zip file itself
		if info.IsDir() || strings.HasSuffix(path, "evidence.zip") {
			return nil
		}

		relPath, err := filepath.Rel(evidenceRoot, path)
		if err != nil {
			return err
		}

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

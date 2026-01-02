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
	base := run.BaseDir()

	// ✅ Timestamped zip name
	zipName := base + ".zip"
	zipPath := filepath.Join(base, zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(run.EvidencePath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and the zip file itself
		if info.IsDir() || strings.HasSuffix(path, zipName) {
			return nil
		}

		relPath, err := filepath.Rel(run.EvidencePath(), path)
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

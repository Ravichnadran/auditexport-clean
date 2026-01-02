package run

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func WriteHashes() error {
	evidenceRoot := EvidencePath()
	hashFilePath := EvidencePath("run", "hashes.txt")

	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		return err
	}
	defer hashFile.Close()

	return filepath.Walk(evidenceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// --------------------------------------------------
		// Skip directories
		// --------------------------------------------------
		if info.IsDir() {
			return nil
		}

		// --------------------------------------------------
		// Skip the hashes file itself
		// --------------------------------------------------
		if strings.HasSuffix(path, "hashes.txt") {
			return nil
		}

		// --------------------------------------------------
		// 🚫 Skip OS junk files (macOS)
		// --------------------------------------------------
		if info.Name() == ".DS_Store" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return err
		}

		hash := hex.EncodeToString(hasher.Sum(nil))

		relPath, err := filepath.Rel(evidenceRoot, path)
		if err != nil {
			return err
		}

		_, err = hashFile.WriteString(hash + "  " + relPath + "\n")
		return err
	})
}

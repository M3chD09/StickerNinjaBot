package filestorage

import (
	"archive/zip"
	"os"
	"path/filepath"
)

func zipDir(dir, zipFilePath string) error {
	// Get a Buffer to Write To
	outFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	if err := addFiles(w, dir, ""); err != nil {
		return err
	}

	// Make sure to check the error on Close.
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		newBase := filepath.Join(basePath, file.Name())
		newBaseInZip := filepath.Join(baseInZip, file.Name())
		if !file.IsDir() {
			dat, err := os.ReadFile(newBase)
			if err != nil {
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(newBaseInZip)
			if err != nil {
				return err
			}
			if _, err := f.Write(dat); err != nil {
				return err
			}
		} else {
			// Recurse
			if err := addFiles(w, newBase, newBaseInZip); err != nil {
				return err
			}
		}
	}

	return nil
}

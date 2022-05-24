package filestorage

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
)

func zipDir(dir, zipFilePath string) {
	// Get a Buffer to Write To
	outFile, err := os.Create(zipFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, dir, "")

	// Make sure to check the error on Close.
	defer w.Close()
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := os.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				log.Fatal(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := filepath.Join(basePath, file.Name())
			addFiles(w, newBase, filepath.Join(baseInZip, file.Name()))
		}
	}
}

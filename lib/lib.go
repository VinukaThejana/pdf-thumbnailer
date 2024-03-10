/*
Copyright © 2024 Vinuka Kodituwakku <vinuka.t@pm.me>
*/

// Package lib contains various libraries
package lib

import (
	"io/fs"
	"path/filepath"
)

// GetAllPDFFiles is a function to retrive an array of all the pdf files from the given location regardless
// of the depth that they are in
func GetAllPDFFiles(root string) []string {
	var files []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == ".pdf" {
			files = append(files, path)
		}

		return nil
	})

	return files
}

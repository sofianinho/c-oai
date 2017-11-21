package utils

import (
	"path/filepath"
	"regexp"
	"os"
)

//FindDirExt finds a directory with a given pattern
func FindDirExt(ext, location string)([]string){
	var files []string
	filepath.Walk(location, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

//FindFileExt finds a file with a given pattern
func FindFileExt(ext, location string)([]string){
	var files []string
	filepath.Walk(location, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}
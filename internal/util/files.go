package util

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// recursively list of all files in a directory
func TraverseDirectory(root *string, extensions *[]string, ignoredPaths *[]string) (*[]string, error) {
	var files []string
	err := filepath.WalkDir(*root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isValidFile(&path, extensions, ignoredPaths) {
			files = append(files, path)
		}
		return nil
	})
	return &files, err
}

// verify if a file is to be analyzed
func isValidFile(path *string, extensions *[]string, ignoredPaths *[]string) bool {
	for _, p := range *ignoredPaths {
		if len(p) == 0 {
			continue
		}
		matched, _ := regexp.Match(p, []byte(*path))
		if matched {
			return false
		}
	}

	for _, ext := range *extensions {
		if strings.HasSuffix(*path, ext) {
			return true
		}
	}
	return false
}

// generate a full path of an imported file
func FullPath(file string, base string, replacers *map[string]string) string {
	if strings.HasPrefix(file, ".") {
		return filepath.Join(filepath.Dir(base), file)
	}
	for key, value := range *replacers {
		if strings.HasPrefix(file, key) {
			file = strings.Replace(file, key, value, 1)
			break
		}
	}
	return file
}

// check whether a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

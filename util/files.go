package util

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Path string
	Code string
}

// recursively list of all files in a directory
func TraverseDirectory(root *string, extensions *[]string, ignoredPaths *[]string) (*[]File, error) {
	var files []File
	err := filepath.WalkDir(*root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isValidFile(&path, extensions, ignoredPaths) {
			code, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, File{path, string(code)})
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

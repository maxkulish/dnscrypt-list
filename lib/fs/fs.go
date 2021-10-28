// Package fs contains all utilities to work with files
package fs

import "os"

// CreateFileOrTruncate creates or truncates the named file
func CreateFileOrTruncate(path string) (*os.File, error) {

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// CreateFileIfNotExist check if the named file exist
// and if not exist creates it
func CreateFileIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		return err
	}

	return nil
}

// MkdirAllIfNotExist creates the given path dir if it isn't exist.
func MkdirAllIfNotExist(path string) error {
	if IsPathExist(path) {
		return nil
	}

	return os.MkdirAll(path, 0755)
}

// IsPathExist returns whether the given path exists.
func IsPathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

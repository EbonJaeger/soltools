package soltools

import (
	"io"
	"os"
	"path/filepath"
)

// LocalRepo is the path to the local Solbuild repo.
const LocalRepo = "/var/lib/solbuild/local"

// CopyPackage copies an eopkg archive to the local Solbuild repo.
func CopyPackage(path string) error {
	_, name := filepath.Split(path)

	src, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer src.Close()

	destPath := filepath.Join(LocalRepo, name)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	err = os.Chmod(destPath, 0644)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dest, src); err != nil {
		return err
	}

	return nil
}

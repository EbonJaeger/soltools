package repo

import (
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// LocalRepo is the path to the local Solbuild repo.
const LocalPath = "/var/lib/solbuild/local"

// Clean removes all package files in the local Solbuild repo.
func Clean() (removed []string, err error) {
	repo := os.DirFS(LocalPath)
	packages, err := fs.Glob(repo, "*.eopkg")
	if err != nil {
		return
	}

	for _, pkg := range packages {
		path := filepath.Join(LocalPath, pkg)
		if err = os.Remove(path); err != nil {
			return
		}
		removed = append(removed, pkg)
	}

	return
}

// CopyInto copies an eopkg archive to the local Solbuild repo.
func CopyInto(path string) error {
	_, name := filepath.Split(path)

	src, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer src.Close()

	destPath := filepath.Join(LocalPath, name)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if err = os.Chmod(destPath, 0644); err != nil {
		return err
	}

	if _, err = io.Copy(dest, src); err != nil {
		return err
	}

	return nil
}

// Index indexes the packages in the local repo.
func Index() error {
	command := exec.Command("eopkg", "index", "--skip-signing", LocalPath)
	command.Dir = LocalPath

	return command.Run()
}

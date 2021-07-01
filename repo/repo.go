package repo

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

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

// Index indexes the packages in the local repo.
func Index() error {
	command := exec.Command("eopkg", "index", "--skip-signing", LocalPath)
	command.Dir = LocalPath

	return command.Run()
}

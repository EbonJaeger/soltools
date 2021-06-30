package soltools

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// CleanRepo removes all package files in the local Solbuild repo.
func CleanRepo() (removed []string, err error) {
	repo := os.DirFS(LocalRepo)
	packages, err := fs.Glob(repo, "*.eopkg")
	if err != nil {
		return
	}

	for _, pkg := range packages {
		path := filepath.Join(LocalRepo, pkg)
		if err = os.Remove(path); err != nil {
			return
		}
		removed = append(removed, pkg)
	}

	return
}

// IndexRepo indexes the packages in the local repo.
func IndexRepo() error {
	command := exec.Command("eopkg", "index", "--skip-signing", LocalRepo)
	command.Dir = LocalRepo

	return command.Run()
}

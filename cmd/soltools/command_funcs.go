package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/soltools/internal/permission"
	"github.com/EbonJaeger/soltools/repo"
	"github.com/go-git/go-git/v5"
)

// CleanPackages removes all packages in the local Solbuild repo and reindexes it.
func CleanPackages(root *cmd.Root, c *cmd.Sub) {
	if err := permission.EscalateIfNeeded(); err != nil {
		log.Fatalf("Unable to escalate privileges: %s\n", err)
	}

	log.Infoln("Looking for packages to clean")

	// Clean the local repo
	cleaned, err := repo.Clean()
	if err != nil {
		log.Errorf("Error cleaning packages: %s\n", err)
	} else {
		log.Goodln("Removed the following packages:")
		log.Printf("\t- %s\n", strings.Join(cleaned, "\n\t- "))
	}

	indexRepo()
}

// CopyPackages copies any eopkg file in the current repo to the local Solbuild
// repo, and indexes the repo after.
func CopyPackages(root *cmd.Root, c *cmd.Sub) {
	if err := permission.EscalateIfNeeded(); err != nil {
		log.Fatalf("Unable to escalate privileges: %s\n", err)
	}

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current directory, %s\n", err)
	}

	log.Infoln("Looking for packages to copy")

	// Glob all .eopkg files
	dir := os.DirFS(cwd)
	packages, err := fs.Glob(dir, "*.eopkg")
	if err != nil {
		log.Fatalf("Error finding packages to copy: %s\n", err)
	}
	if packages == nil {
		log.Infoln("No packages to copy!")
		return
	}

	log.Goodln("Found the following packages:")
	log.Printf("\t- %s\n", strings.Join(packages, "\n\t- "))

	// Copy each package to the local repo
	for _, eopkg := range packages {
		err = repo.CopyInto(filepath.Join(cwd, eopkg))
		if err != nil {
			log.Errorf("Error copying package '%s': %s\n", eopkg, err)
			continue
		}
	}

	indexRepo()
}

func indexRepo() {
	log.Infoln("Indexing local repo")
	if err := repo.Index(); err != nil {
		log.Fatalf("Error indexing local repo: %s\n", err)
	} else {
		log.Goodln("Local repo indexed")
	}
}

// InitRepo creates and initializes a new package repository.
func InitRepo(root *cmd.Root, c *cmd.Sub) {
	flags := c.Flags.(*InitFlags)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current directory: %s\n", err)
	}

	if _, err = os.Stat(filepath.Join(cwd, "common")); err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("common not found, aborting clone")
		} else {
			log.Fatalf("Error trying to find common directory: %s\n", err)
		}
	}

	name := c.Args.(*InitArgs).Name
	path := filepath.Join(cwd, name)

	log.Infoln("Creating git repo")

	if _, err = git.PlainInit(path, false); err != nil {
		log.Fatalf("Error creating new git repo: %s\n", err)
	}
	log.Goodln("Git repo created")

	log.Infoln("Creating Makefile in the repo")
	makefile, err := os.Create(filepath.Join(path, "Makefile"))
	if err != nil {
		log.Fatalf("Error creating Makefile: %s\n", err)
	}
	defer makefile.Close()

	if _, err = makefile.WriteString(makefileString); err != nil {
		log.Fatalf("Error writing to Makefile: %s\n", err)
	}
	log.Goodln("Makefile written")

	// Write maintainers file unless told not to
	if !flags.SkipMaintainers {
		log.Infoln("Creating maintainers file")
		maintainersFile, err := os.Create(filepath.Join(path, "MAINTAINERS.md"))
		if err != nil {
			log.Fatalf("Error creating maintainers file: %s\n", err)
		}
		defer maintainersFile.Close()

		if _, err = maintainersFile.WriteString(maintainerString); err != nil {
			log.Fatalf("Error writing to maintainer file: %s\n", err)
		}
		log.Goodln("Maintainers file written")
	} else {
		log.Infoln("Skipping writing maintainers file")
	}

	log.Infoln("Running yauto.py to generate package.yml")
	cmd := exec.Command(filepath.Join(cwd, "common", "Scripts", "yauto.py"), c.Args.(*InitArgs).URL)
	cmd.Dir = path
	if err = cmd.Run(); err != nil {
		log.Fatalf("Error generating package.yml: %s", err)
	} else {
		log.Goodln("Package repo initialized")
	}
}

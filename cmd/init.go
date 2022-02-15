package cmd

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/go-git/go-git/v5"
)

var (
	//go:embed embed/Makefile
	makefileString string

	//go:embed embed/MAINTAINERS.md
	maintainerString string
)

// InitArgs holds the arguments for the init command.
type InitArgs struct {
	Name string `desc:"Name of the package to clone"`
	URL  string `desc:"URL of the source tarball to use"`
}

type InitFlags struct {
	SkipMaintainers bool `long:"skip-maintainers" desc:"Skip creating a MAINTAINERS.md file for this repository"`
}

// Init is our command to initialize a new local package repository.
var Init = cmd.Sub{
	Name:  "init",
	Alias: "i",
	Short: "Initializes a new package repo",
	Args:  &InitArgs{},
	Flags: &InitFlags{},
	Run:   InitRepo,
}

// InitRepo creates and initializes a new package repository.
func InitRepo(root *cmd.Root, c *cmd.Sub) {
	logger := NewLogger()
	flags := c.Flags.(*InitFlags)

	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Unable to get current directory: %s\n", err)
	}

	if _, err = os.Stat(filepath.Join(cwd, "common")); err != nil {
		if os.IsNotExist(err) {
			logger.Fatalln("common not found, aborting clone")
		} else {
			logger.Fatalf("Error trying to find common directory: %s\n", err)
		}
	}

	name := c.Args.(*InitArgs).Name
	path := filepath.Join(cwd, name)

	logger.Infoln("Creating git repo")

	if _, err = git.PlainInit(path, false); err != nil {
		logger.Fatalf("Error creating new git repo: %s\n", err)
	}
	logger.Goodln("Git repo created")

	logger.Infoln("Creating Makefile in the repo")
	makefile, err := os.Create(filepath.Join(path, "Makefile"))
	if err != nil {
		logger.Fatalf("Error creating Makefile: %s\n", err)
	}
	defer makefile.Close()

	if _, err = makefile.WriteString(makefileString); err != nil {
		logger.Fatalf("Error writing to Makefile: %s\n", err)
	}
	logger.Goodln("Makefile written")

	// Write maintainers file unless told not to
	if !flags.SkipMaintainers {
		logger.Infoln("Creating maintainers file")
		maintainersFile, err := os.Create(filepath.Join(path, "MAINTAINERS.md"))
		if err != nil {
			logger.Fatalf("Error creating maintainers file: %s\n", err)
		}
		defer maintainersFile.Close()

		if _, err = maintainersFile.WriteString(maintainerString); err != nil {
			logger.Fatalf("Error writing to maintainer file: %s\n", err)
		}
		logger.Goodln("Maintainers file written")
	} else {
		logger.Infoln("Skipping writing maintainers file")
	}

	logger.Infoln("Running yauto.py to generate package.yml")
	cmd := exec.Command(filepath.Join(cwd, "common", "Scripts", "yauto.py"), c.Args.(*InitArgs).URL)
	cmd.Dir = path
	if err = cmd.Run(); err != nil {
		logger.Fatalf("Error generating package.yml: %s", err)
	} else {
		logger.Goodln("Package repo initialized")
	}
}

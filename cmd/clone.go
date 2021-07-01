package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/soltools"
	"github.com/go-git/go-git/v5"
)

// CloneArgs holds the arguments for the clone command.
type CloneArgs struct {
	Name string `desc:"Name of the package to clone"`
}

// Clone is our repo clone command.
var Clone = cmd.Sub{
	Name:  "clone",
	Short: "Clone a package from the official Solus repo",
	Args:  &CloneArgs{},
	Run:   ClonePackage,
}

// ClonePackage clones a package repository from the upstream Solus repo.
func ClonePackage(root *cmd.Root, c *cmd.Sub) {
	logger := soltools.NewLogger()

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

	name := c.Args.(*CloneArgs).Name
	path := filepath.Join(cwd, name)

	logger.Infoln("Cloning Solus package repository")

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://dev.getsol.us/source/%s.git", name),
		Progress: os.Stdout,
	})

	if err != nil {
		logger.Fatalf("Error cloning repository: %s\n", err)
	} else {
		logger.Goodln("Package repository cloned")
	}
}

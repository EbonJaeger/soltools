package cmd

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/soltools"
)

// Copy is the copy subcommand.
var Copy = cmd.Sub{
	Name:  "copy",
	Alias: "c",
	Short: "Copies .eopkg files to the local repo and indexes",
	Run:   CopyPackages,
}

// CopyPackages copies any eopkg file in the current repo to the local Solbuild
// repo, and indexes the repo after.
func CopyPackages(root *cmd.Root, c *cmd.Sub) {
	logger := soltools.NewLogger()

	user, err := user.Current()
	if err != nil {
		logger.Fatalf("Unable to check for root privileges: %s\n", err)
	}
	if user.Name != "root" {
		logger.Fatalln("This command must be run with elevated privileges")
	}

	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Unable to get current directory, %s\n", err)
	}

	logger.Infoln("Looking for packages to copy")

	dir := os.DirFS(cwd)
	packages, err := fs.Glob(dir, "*.eopkg")
	if err != nil {
		logger.Fatalf("Error finding packages to copy: %s\n", err)
	}
	if packages == nil {
		logger.Infoln("No packages to copy!")
		return
	}

	logger.Goodln("Found the following packages:")
	logger.Printf("\t- %s\n", strings.Join(packages, "\n\t- "))

	for _, eopkg := range packages {
		err = soltools.CopyPackage(filepath.Join(cwd, eopkg))
		if err != nil {
			logger.Errorf("Error copying package '%s': %s\n", eopkg, err)
			continue
		}
	}

	logger.Infoln("Indexing local repo")
	err = soltools.IndexRepo()
	if err != nil {
		logger.Fatalf("Error indexing local repo: %s\n", err)
	} else {
		logger.Goodln("Indexed local repo")
	}
}

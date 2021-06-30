package cmd

import (
	"os/user"
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/soltools"
)

// Clean is the clean subcommand.
var Clean = cmd.Sub{
	Name:  "clean",
	Short: "Removes .eopkg files from the local repo and re-indexes",
	Run:   CleanPackages,
}

// CleanPackages removes all packages in the local Solbuild repo and reindexes it.
func CleanPackages(root *cmd.Root, c *cmd.Sub) {
	logger := soltools.NewLogger()

	user, err := user.Current()
	if err != nil {
		logger.Fatalf("Unable to check for root privileges: %s\n", err)
	}
	if user.Name != "root" {
		logger.Fatalln("This command must be run with elevated privileges")
	}

	logger.Infoln("Looking for packages to clean")

	cleaned, err := soltools.CleanRepo()
	if err != nil {
		logger.Errorf("Error cleaning packages: %s\n", err)
	} else {
		logger.Goodln("Removed the following packages:")
		logger.Printf("\t- %s\n", strings.Join(cleaned, "\n\t- "))
	}

	logger.Infoln("Indexing local repo")
	err = soltools.IndexRepo()
	if err != nil {
		logger.Fatalf("Error indexing local repo: %s\n", err)
	} else {
		logger.Goodln("Indexed local repo")
	}
}

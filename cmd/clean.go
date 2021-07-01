package cmd

import (
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/soltools/internal/permission"
	"github.com/EbonJaeger/soltools/repo"
)

// Clean is the clean subcommand.
var Clean = cmd.Sub{
	Name:  "clean",
	Short: "Removes .eopkg files from the local repo and re-indexes",
	Run:   CleanPackages,
}

// CleanPackages removes all packages in the local Solbuild repo and reindexes it.
func CleanPackages(root *cmd.Root, c *cmd.Sub) {
	logger := NewLogger()

	if err := permission.EscalateIfNeeded(); err != nil {
		logger.Fatalf("Unable to escalate privileges: %s\n", err)
	}

	logger.Infoln("Looking for packages to clean")

	cleaned, err := repo.Clean()
	if err != nil {
		logger.Errorf("Error cleaning packages: %s\n", err)
	} else {
		logger.Goodln("Removed the following packages:")
		logger.Printf("\t- %s\n", strings.Join(cleaned, "\n\t- "))
	}

	logger.Infoln("Indexing local repo")
	if err = repo.Index(); err != nil {
		logger.Fatalf("Error indexing local repo: %s\n", err)
	} else {
		logger.Goodln("Local repo indexed")
	}
}

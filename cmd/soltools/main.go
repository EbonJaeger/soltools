package main

import (
	_ "embed"
	log2 "log"
	"os"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
)

var (
	//go:embed embed/Makefile
	makefileString string

	//go:embed embed/MAINTAINERS.md
	maintainerString string

	log *waterlog.WaterLog
)

// InitArgs holds the arguments for the init command.
type InitArgs struct {
	Name string `desc:"Name of the package to clone"`
	URL  string `desc:"URL of the source tarball to use"`
}

type InitFlags struct {
	SkipMaintainers bool `long:"skip-maintainers" desc:"Skip creating a MAINTAINERS.md file for this repository"`
}

func init() {
	// Set up the loggers
	log = waterlog.New(os.Stdout, "soltools", log2.Ltime)
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)
}

func main() {
	root := &cmd.Root{
		Name:  "soltools",
		Short: "Tool to assist with Solus packaging",
	}

	cmd.Register(&cmd.Sub{
		Name:  "clean",
		Short: "Removes .eopkg files from the local repo and re-indexes",
		Run:   CleanPackages,
	})

	cmd.Register(&cmd.Sub{
		Name:  "copy",
		Alias: "c",
		Short: "Copies .eopkg files to the local repo and re-indexes",
		Run:   CopyPackages,
	})

	cmd.Register(&cmd.Sub{
		Name:  "init",
		Alias: "i",
		Short: "Initializes a new package repo",
		Args:  &InitArgs{},
		Flags: &InitFlags{},
		Run:   InitRepo,
	})

	root.Run()
}

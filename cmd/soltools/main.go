package main

import (
	"github.com/DataDrake/cli-ng/v2/cmd"
	commands "github.com/EbonJaeger/soltools/cmd"
)

type GlobalFlags struct {
}

func main() {
	root := &cmd.Root{
		Name:  "soltools",
		Short: "Tool to assist with Solus packaging",
		Flags: &GlobalFlags{},
	}

	cmd.Register(&commands.Copy)
	cmd.Register(&commands.Clean)
	cmd.Register(&commands.Clone)
	cmd.Register(&commands.Init)

	root.Run()
}

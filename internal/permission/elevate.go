package permission

import (
	"os"
	"os/exec"
)

// EscalateIfNeeded checks if the program is running as root, and re-runs
// the program with sudo if it is not.
//
// If the program is already running with root privileges, the program will
// continue as normal.
func EscalateIfNeeded() error {
	if os.Getuid() == 0 && os.Geteuid() == 0 {
		return nil
	}

	cmd := exec.Command("/usr/bin/sudo")

	path, err := os.Executable()
	if err != nil {
		return err
	}

	cmd.Args = append(cmd.Args, path)
	cmd.Args = append(cmd.Args, os.Args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

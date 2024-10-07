package command

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Returns the shell and shell flag to use for the current execution environment
func GetShell() (string, string) {
	if runtime.GOOS == "windows" {
		// Windows: prefer PowerShell if available, else fallback to cmd.exe
		if os.Getenv("ComSpec") != "" {
			return "cmd.exe", "/C"
		}
		return "powershell.exe", "-Command"
	} else {
		// Unix-like systems: use the SHELL environment variable or fallback to /bin/sh
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh" // fallback to /bin/sh if SHELL is not set
		}

		return shell, "-c"
	}
}

// runCommand runs the generated command in the detected shell
func RunCommand(command string) error {
	shell, shellFlag := GetShell()

	// Create the command to execute using the detected shell
	cmd := exec.Command(shell, shellFlag, command)

	// Set up to pipe the output and error directly to stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

package commander

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type ShellError struct {
	Err     error
	Command string
}

func (e *ShellError) Error() string {
	return fmt.Sprintf("Cannot run %q. Error %v", e.Command, e.Err)
}

type Commander interface {
	Exec(cmd *exec.Cmd) (string, error)
	ExecSilently(cmd *exec.Cmd) error
}

type DefaultCommander struct {
	logger *log.Logger
}

func (c DefaultCommander) Exec(cmd *exec.Cmd) (string, error) {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if c.logger != nil {
			c.logger.Println(err, string(output))
		}
		// return "", &ShellError{strings.Join(cmd.Args, " "), err}
		return "", &ShellError{Command: strings.Join(cmd.Args, " "), Err: err}
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

func (c DefaultCommander) ExecSilently(cmd *exec.Cmd) error {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, " "))
	}

	if err := cmd.Run(); err != nil {
		if c.logger != nil {
			c.logger.Println(err)
		}
		// return &ShellError{strings.Join(cmd.Args, " "), Err: err}
		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}
	return nil
}

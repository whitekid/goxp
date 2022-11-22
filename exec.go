package goxp

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/whitekid/goxp/log"
)

var (
	loggerExec = log.New(log.AddCallerSkip(1))
)

func Exec(command ...string) *Executor { return &Executor{command: command} }

// Executor just simple command executor
type Executor struct {
	shell   bool
	command []string
	dir     string
}

func (exc *Executor) Shell() *Executor         { exc.shell = true; return exc }
func (exc *Executor) NoShell() *Executor       { exc.shell = false; return exc }
func (exc *Executor) Dir(dir string) *Executor { exc.dir = dir; return exc }

func (exc *Executor) buildCmd(ctx context.Context) *exec.Cmd {
	var name string
	var args []string

	if exc.shell {
		name = "sh"
		args = append([]string{"-c"}, exc.command...)
	} else if len(exc.command) > 0 {
		name = exc.command[0]

		if len(exc.command) > 1 {
			args = exc.command[1:]
		} else {
			args = nil
		}
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = exc.dir

	return cmd
}

// Do execute command, output stdout to stdout and stderr to stderr
func (exc *Executor) Do(ctx context.Context) error {
	dir := exc.dir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	loggerExec.Debugf("execute: %s", strings.Join(exc.command, " "))
	loggerExec.Debugf("dir: %s", dir)

	cmd := exc.buildCmd(ctx)

	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Output run command and return stdout
func (exc *Executor) Output(ctx context.Context) ([]byte, error) {
	dir := exc.dir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	loggerExec.Debugf("execute: %s", strings.Join(exc.command, " "))
	loggerExec.Debugf("dir: %s", dir)

	return exc.buildCmd(ctx).Output()
}

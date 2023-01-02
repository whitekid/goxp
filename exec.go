package goxp

import (
	"context"
	"io"
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

	// pipe
	stdin          func(io.WriteCloser)
	stdout, stderr func(io.ReadCloser)
}

func (exc *Executor) Shell(shell bool) *Executor { exc.shell = shell; return exc }
func (exc *Executor) Dir(dir string) *Executor   { exc.dir = dir; return exc }

func (exc *Executor) Pipe(stdin func(io.WriteCloser), stdout, stderr func(io.ReadCloser)) *Executor {
	exc.stdin = stdin
	exc.stdout = stdout
	exc.stderr = stderr
	return exc
}

func (exc *Executor) buildCmd(ctx context.Context) (*exec.Cmd, error) {
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
	loggerExec.Debugf("execute: %s", strings.Join(exc.command, " "))

	cmd.Dir = exc.dir
	dir := cmd.Dir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	loggerExec.Debugf("dir: %s", dir)

	if exc.stdin != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() { exc.stdin(stdin) }()
	}

	if exc.stdout != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		go func() { exc.stdout(stdout) }()
	}

	if exc.stderr != nil {
		stdout, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}

		go func() { exc.stderr(stdout) }()
	}

	return cmd, nil
}

// Do execute command
func (exc *Executor) Do(ctx context.Context) error {
	cmd, err := exc.buildCmd(ctx)
	if err != nil {
		return err
	}

	if exc.stdout == nil {
		cmd.Stderr = os.Stdout
	}

	if exc.stderr == nil {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Output shortcut for Cmd.Output()
func (exc *Executor) Output(ctx context.Context) ([]byte, error) {
	cmd, err := exc.buildCmd(ctx)
	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

// CombinedOutput shortcut for Cmd.CombinedOutput()
func (exc *Executor) CombinedOutput(ctx context.Context) ([]byte, error) {
	cmd, err := exc.buildCmd(ctx)
	if err != nil {
		return nil, err
	}

	return cmd.CombinedOutput()
}

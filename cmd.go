package mach

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"unicode"
)

type cmd struct {
	cmd  string
	args []string
}

func (c cmd) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	if len(c.cmd) == 0 {
		return errors.New("command must not be empty")
	}

	var errBuffer bytes.Buffer

	cmd := exec.CommandContext(ctx, c.cmd, c.args...)

	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		if errBuffer.Len() != 0 {
			return fmt.Errorf("%s: %w", strings.TrimRightFunc(errBuffer.String(), unicode.IsSpace), err)
		} else {
			return err
		}
	}

	if errBuffer.Len() != 0 {
		return errors.New(errBuffer.String())
	}

	return nil
}

package mach

import (
	"bytes"
	"context"
	"errors"
	"io"
)

type Filter func([]byte) (bool, error)

func (f Filter) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	return ByLine(in, func(b []byte) error {
		ok, err := f(b)
		if err != nil {
			return err
		}

		if ok {
			_, err := out.Write(b)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// The slice passed to fn is only valid for the call of the function!
func ByLine(in io.Reader, fn func([]byte) error) error {
	buf := make([]byte, 256)
	var line bytes.Buffer

	for {
		n, err := in.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line.Len() != 0 {
					return fn(line.Bytes())
				}
				break
			}

			return err
		}

		lastLine := 0

		for i := 0; i <= n; i++ {
			if buf[i] == '\n' {
				line.Write(buf[lastLine:i])
				if err := fn(line.Bytes()); err != nil {
					return err
				}
				line.Reset()
				lastLine = i
			}
		}

		if lastLine < n {
			line.Write(buf[lastLine:n])
		}
	}

	return nil
}

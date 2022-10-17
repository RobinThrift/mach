package mach

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
)

type Runnable interface {
	// Run the command with the provided in and out streams. Errors will be written to out.
	Run(ctx context.Context, in io.Reader, out io.Writer) error
}

type RunnableFunc func(ctx context.Context, in io.Reader, out io.Writer) error

func (r RunnableFunc) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	return r(ctx, in, out)
}

type Runner struct {
	next     *Runner
	runnable Runnable
}

func NewRunner(rr Runnable) *Runner {
	return &Runner{
		runnable: rr,
	}
}

// Pipe the result of the previous step into this step.
func (r *Runner) Pipe(rr Runnable) *Runner {
	for n := r; n != nil; n = r.next {
		if n.next == nil {
			if rn, ok := rr.(*Runner); ok {
				r.next = rn
				return r
			} else {
				n.next = NewRunner(rr)
			}
		}
	}

	return r
}

// Go executes the pipeline using os.Stdin and os.Stdout. Errors will be written to out.
func (r *Runner) Go() {
	r.GoCtx(context.Background())
}

func (r *Runner) GoCtx(ctx context.Context) {
	err := r.Run(ctx, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RunBytes executes the pipeline with the provided in streams and return a byte slice with the result. Any error will be returned as the error of this call.
func (r *Runner) RunBytes(ctx context.Context, in io.Reader) ([]byte, error) {
	var out bytes.Buffer
	err := r.Run(ctx, in, &out)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// RunString exexutes the pipeline with the provided in streams and return a string with the result. Any error will be returned as the error of this call.
func (r *Runner) RunString(ctx context.Context, in io.Reader) (string, error) {
	res, err := r.RunBytes(ctx, in)
	return string(res), err
}

func (r *Runner) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	nextIn := in
	for rr := r; rr != nil; rr = r.next {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		default:
		}

		if rr.next == nil {
			if err := rr.run(ctx, nextIn, out); err != nil {
				select {
				case err := <-errChan:
					return err
				default:
				}

				return err
			}
			break
		}

		pipeReader, pipeWriter := io.Pipe()
		defer pipeReader.Close()
		go func(ctx context.Context, rr *Runner, in io.Reader, out io.WriteCloser) {
			defer out.Close()
			if err := rr.run(ctx, in, out); err != nil {
				cancel()
				out.Close()
				errChan <- err
			}
		}(ctx, rr, nextIn, pipeWriter)
		nextIn = pipeReader
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	default:
	}

	return nil
}

func (r *Runner) run(ctx context.Context, in io.Reader, out io.Writer) error {
	return r.runnable.Run(ctx, in, out)
}

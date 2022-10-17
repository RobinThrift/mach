package mach

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunner_SingleCommand(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()

	var in bytes.Buffer
	var out bytes.Buffer

	err := Run("ls ./examples").Run(ctx, &in, &out)
	assert.NoError(t, err)

	assert.Equal(t, `simple
`, out.String())
}

func TestRunner_SimplePipeline(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()
	var in bytes.Buffer
	var out bytes.Buffer

	err := Run("ls examples").Pipe(Run("wc -l")).Run(ctx, &in, &out)
	assert.NoError(t, err)

	assert.Equal(t, `       1
`, out.String())
}

func TestRunner_SimplePipeline_WithError(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()
	var in bytes.Buffer
	var out bytes.Buffer

	err := Run("ls examples").Pipe(Run("nonexistetcommandhopefully")).Run(ctx, &in, &out)
	assert.Error(t, err)
}

func createTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	if deadline, ok := t.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}

	return context.WithTimeout(context.Background(), 10*time.Second)
}

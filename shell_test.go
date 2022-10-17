package mach

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()

	r := Shell("ls examples | wc -l")

	var in bytes.Buffer
	var out bytes.Buffer

	err := r.Run(ctx, &in, &out)
	assert.NoError(t, err)

	assert.Equal(t, `       1
`, out.String())
}

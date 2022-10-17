package mach

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()

	r := Run("ls examples")

	var in bytes.Buffer
	var out bytes.Buffer

	err := r.Run(ctx, &in, &out)
	assert.NoError(t, err)

	assert.Equal(t, `simple
`, out.String())
}

func TestQuoteAwareSplit(t *testing.T) {
	tt := []struct {
		in  string
		exp []string
	}{
		{"a b c d e f", []string{"a", "b", "c", "d", "e", "f"}},
		{`a "b c" d e f`, []string{"a", "b c", "d", "e", "f"}},
		{"a b 'c d' e f", []string{"a", "b", "c d", "e", "f"}},
		{`a b 'c "test" d' e f`, []string{"a", "b", `c "test" d`, "e", "f"}},
		{"a b 'c d e f", []string{"a", "b", "c d e f"}},
		{`a b 'c d e f"`, []string{"a", "b", `c d e f"`}},
		{`a b 'c d' e f"`, []string{"a", "b", "c d", "e", `f"`}},
	}

	for _, tt := range tt {
		actual := quoteAwareSplit(tt.in)
		assert.Equal(t, tt.exp, actual)
	}
}

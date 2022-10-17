package mach

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	ctx, cancel := createTestContext(t)
	defer cancel()

	var in bytes.Buffer
	var out bytes.Buffer

	in.WriteString("this is a test line\n")
	in.WriteString("this is another test line\n")
	in.WriteString("this line ")
	in.WriteString("is split\n")
	in.WriteString("this is \ntwo lines")
	in.WriteString("\nthe end")

	err := Filter(func(b []byte) (bool, error) {
		return len(b) > 10, nil
	}).Run(ctx, &in, &out)

	assert.NoError(t, err)

	expected := `this is a test line
this is another test line
this line is split`

	assert.Equal(t, expected, out.String())
}

func TestByLine(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer

	in.WriteString("this is a test line\n")
	in.WriteString("this is another test line\n")
	in.WriteString("this line ")
	in.WriteString("is split\n")
	in.WriteString("this is \ntwo lines")
	in.WriteString("\nthe end")

	expected := in.String()

	err := ByLine(&in, func(b []byte) error {
		_, err := out.Write(b)
		return err
	})

	assert.NoError(t, err)

	assert.Equal(t, expected, out.String())
}

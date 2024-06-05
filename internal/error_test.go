package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	errCommon := Error{kind: ENOTFOUND, message: "im good", code: "zero"}
	var err error = &errCommon
	kind, message, code := ParseError(err)
	assert.Equal(t, errCommon.kind, kind)
	assert.Equal(t, errCommon.message, message)
	assert.Equal(t, errCommon.code, code)

	kind, message, code = ParseError(os.ErrPermission)
	assert.Equal(t, DefaultErrorKind, kind)
	assert.Equal(t, DefaultErrorMessage, message)
	assert.Equal(t, DefaultErrorCode, code)
}

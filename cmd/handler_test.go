package cmd

import (
	"testing"

	"github.com/chris-greaves/stencil/cmd/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	mockEngine = new(mocks.Engine)
	mockConfig = new(mocks.Config)
)

func TestNewRootHandlerReturnsValidRootHandler(t *testing.T) {
	handler := NewRootHandler(mockConfig, mockEngine)

	assert.NotNil(t, handler, "returned handler should not be nil")
	assert.IsType(t, RootHandler{}, handler)
}

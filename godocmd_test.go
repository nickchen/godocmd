package godocmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixture(t *testing.T) {

	d, err := New("./markdown.tmpl")
	assert.NoError(t, err, "failed to initialize")
	err = d.ProcessPackageDirs("./docs/", "github.com/nickchen/godocmd", "./fixture")
	assert.NoError(t, err, "failed to process")
	assert.NotNil(t, nil, "failed")
}

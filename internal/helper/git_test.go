package helper_test

import (
	"testing"

	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestMakeSSHRepo(t *testing.T) {
	url := "http://github.com/username/repo.git"
	expectedUrl := "git@github.com:username/repo.git"
	actualUrl := helper.MakeSSHRepo(url)
	assert.Equal(t, expectedUrl, actualUrl, "should work")
}

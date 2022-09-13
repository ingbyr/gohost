package gohost

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMasks(t *testing.T) {
	a := assert.New(t)
	node := NewTreeNode(&Group{})
	a.False(node.IsEnabled())
	node.SetEnabled(true)
	a.True(node.IsEnabled())

	a.False(node.IsFolded())
	node.SetFolded(true)
	a.True(node.IsFolded())
	node.SetFolded(false)
	a.False(node.IsFolded())
}

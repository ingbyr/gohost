package gohost

import (
	"github.com/stretchr/testify/assert"
	"github.com/timshannon/bolthold"
	"testing"
)

func TestService_SaveHost(t *testing.T) {
	a := assert.New(t)
	GetService().store.DeleteMatching(&LocalHost{}, &bolthold.Query{})
	hosts := []Host{
		&LocalHost{
			ID:      1000,
			Name:    "host-1000",
			Content: []byte("127.0.0.1 localhost"),
			Desc:    "host1000",
			GroupID: 4,
		},
		&LocalHost{
			ID:      1001,
			Name:    "host-1",
			Content: []byte("127.0.0.2 localhost"),
			Desc:    "host1001",
			GroupID: 0,
		},
	}

	for _, host := range hosts {
		if err := GetService().SaveNode(NewTreeNode(host)); err != nil {
			a.NoError(err)
		}
	}
}

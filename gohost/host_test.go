package gohost

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/timshannon/bolthold"
	"testing"
)

func TestService_SaveHost(t *testing.T) {
	a := assert.New(t)
	GetService().store.DeleteMatching(&LocalHost{}, &bolthold.Query{})
	hosts := []Host{
		&LocalHost{
			ID:      "1000",
			Name:    "host-1000",
			Content: []byte("127.0.0.1 localhost"),
			Desc:    "host1000",
			GroupID: "3",
		},
	}

	for _, host := range hosts {
		if err := GetService().SaveHost(host); err != nil {
			a.NoError(err)
		}
	}
}

func TestService_LoadHost(t *testing.T) {
	hosts := GetService().loadLocalHosts("3")
	for _, host := range hosts {
		fmt.Printf("%+v\n", host)
	}
}

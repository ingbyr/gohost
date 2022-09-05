package store

import (
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"gohost/config"
	"os"
	"sync"
)

var (
	instance *store
	once     sync.Once
)

func Store() *store {
	once.Do(func() {
		cfg := config.Config()
		instance = New(&options{
			File: cfg.DBFile,
			Options: &bolthold.Options{
				Encoder: nil,
				Decoder: nil,
				Options: &bbolt.Options{},
			},
		})
	})
	return instance
}

type options struct {
	File string
	*bolthold.Options
}

type store struct {
	*bolthold.Store
}

func New(opt *options) *store {
	s, err := bolthold.Open(opt.File, os.ModePerm, opt.Options)
	if err != nil {
		panic(err)
	}
	return &store{Store: s}
}

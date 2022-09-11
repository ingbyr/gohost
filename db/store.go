package db

import (
	"errors"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"gohost/config"
	"os"
	"sync"
)

type ID = uint64

var (
	instance *Store
	once     sync.Once
)

func Instance() *Store {
	once.Do(func() {
		cfg := config.Instance()
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

type Store struct {
	*bolthold.Store
}

func New(opt *options) *Store {
	s, err := bolthold.Open(opt.File, os.ModePerm, opt.Options)
	if err != nil {
		panic(err)
	}
	return &Store{Store: s}
}

func (s *Store) FindNullable(result interface{}, query *bolthold.Query) error {
	err := s.Find(result, query)

	if !errors.Is(bolthold.ErrNotFound, err) {
		return err
	}
	return nil
}

func (s *Store) NextID() any {
	return bolthold.NextSequence()
}

package main

import (
	"github.com/timshannon/bolthold"
	"os"
)

type StoreOptions struct {
	File string
	*bolthold.Options
}

func NewStore(opt *StoreOptions) (*bolthold.Store, error) {
	return bolthold.Open(opt.File, os.ModePerm, opt.Options)
}

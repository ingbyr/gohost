package dal

import (
	"github.com/timshannon/bolthold"
	"os"
)

type Options struct {
	File string
	*bolthold.Options
}

func New(opt *Options) (*bolthold.Store, error) {
	return bolthold.Open(opt.File, os.ModePerm, opt.Options)
}

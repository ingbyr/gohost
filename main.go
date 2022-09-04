package main

import (
	"github.com/timshannon/bolthold"
)

func main() {
	//cmd.Execute()
	store, err := NewStore(&StoreOptions{
		File:    cfg.DBFile,
		Options: &bolthold.Options{},
	})
	defer store.Close()
	if err != nil {
		panic(err)
	}
}

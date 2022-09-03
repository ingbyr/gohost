package main

import (
	"github.com/timshannon/bolthold"
	"gohost/dal"
)

func main() {
	//cmd.Execute()
	store, err := dal.New(&dal.Options{
		File:    cfg.DBFile,
		Options: &bolthold.Options{},
	})
	defer store.Close()
	if err != nil {
		panic(err)
	}
}

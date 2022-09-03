package dal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"io/fs"
	"os"
	"testing"
	"time"
)

type Item struct {
	ID       int
	Category string `boltholdIndex:"Category"`
	Created  time.Time
}

func TestBoltHold(t *testing.T) {
	a := assert.New(t)
	data := []Item{
		{
			ID:       0,
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		{
			ID:       1,
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		{
			ID:       2,
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		{
			Category: "blue",
			ID:       3,
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}
	bho := &bolthold.Options{
		Options: &bbolt.Options{},
	}

	store, err := bolthold.Open("test.db", fs.ModePerm, bho)
	defer func() {
		store.Close()
		os.Remove("test.db")
	}()
	a.NoError(err)

	err = store.Bolt().Update(func(tx *bbolt.Tx) error {
		for i := range data {
			err = store.TxInsert(tx, data[i].ID, data[i])
			a.NoError(err)
		}
		return nil
	})
	a.NoError(err)

	//var res []Item
	//err = store.Find(&res, bolthold.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1*time.Hour)))
	//a.NoError(err)
	//
	//for _, r := range res {
	//	fmt.Println(r)
	//}

	//res = []Item{}
	//err = store.Find(&res, bolthold.Where(bolthold.Key).Ne(2))
	//a.NoError(err)
	//for _, r := range res {
	//	fmt.Println(r)
	//}
	//
	//err = store.Insert(bolthold.NextSequence(), "!23")
	//a.NoError(err)

	err = store.ForEach(bolthold.Where("Category").Eq("blue"), func(record *Item) error{
		fmt.Println(record)
		return nil
	})
	a.NoError(err)
}

package biz

import (
	"context"
	"os"
	"testing"

	"github.com/ingbyr/gohost/dal"
	"github.com/ingbyr/gohost/dal/model"
	"github.com/ingbyr/gohost/dal/query"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.Remove("../test.db")
	if err != nil {
		panic(err)
	}
	dal.DB = dal.ConnectDB("../test.db")
	dal.AutoMigrate()
	query.SetDefault(dal.DB)
}

func TestHost_TableName(t *testing.T) {
	as := assert.New(t)
	ctx := context.Background()
	var err error

	h := query.Host

	_, err = h.Groups.Model(&model.Host{Name: "not-exist"}).Find()
	as.NoError(err)

	h1 := &model.Host{
		Name:    "h1",
		Content: "127.0.0.1 localhost",
		Groups:  []*model.Group{},
	}

	err = h.WithContext(ctx).Create(h1)
	as.NoError(err)

	g1 := &model.Group{
		Name:  "g1",
		Hosts: []*model.Host{},
	}

	g2 := &model.Group{
		Name:  "g2",
		Hosts: []*model.Host{},
	}
	err = h.Groups.Model(h1).Append(g1, g2)
	as.NoError(err)

	hs, err := h.WithContext(ctx).Preload(h.Groups).Where(h.Name.Eq("h1")).Find()
	if err != nil {
		t.Fatal(err)
	}
	as.NotNil(hs)
}

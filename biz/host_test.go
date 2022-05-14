package biz

import (
	"fmt"
	"testing"

	"github.com/ingbyr/gohost/dal"
	"github.com/ingbyr/gohost/dal/model"
	"github.com/ingbyr/gohost/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB("../test.db")
	dal.AutoMigrate()
	query.SetDefault(dal.DB)
}

func TestHost_TableName(t *testing.T) {
	h := query.Host
	res, err := h.Groups.Model(&model.Host{Name: "host-1"}).Find()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(res)
}

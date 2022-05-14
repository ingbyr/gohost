/*
 @Author: ingbyr
*/

package main

import (
	"github.com/ingbyr/gohost/dal"
	"github.com/ingbyr/gohost/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB("test.db").Debug()
	dal.AutoMigrate()
	query.SetDefault(dal.DB)
}

func main() {
	// cmd.Execute()
}

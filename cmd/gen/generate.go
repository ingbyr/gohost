package main

import (
	"github.com/ingbyr/gohost/dal/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery,
	})
	g.ApplyBasic(model.Product{})
	g.Execute()
}

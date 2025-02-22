package main

import (
	"github.com/xLeSHka/calc/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.App).Run()
}

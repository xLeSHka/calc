package main

import (
	"github.com/xLeSHka/calc/app/orchestrator"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		orchestrator.Orchestrator,
	).Run()
}

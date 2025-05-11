package main

import (
	"github.com/xLeSHka/calc/internal/app/orchestrator"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		orchestrator.Orchestrator,
	).Run()
}
